package game

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"flag"
	"strconv"
)

var port int
func init () {
	flag.IntVar(&port, "port", 80, "Port to listen on")
	flag.Parse()
}

type Gameplay func([]*Player, chan Command, chan bool)

type Game struct {
	Playfunc   Gameplay
	NumPlayers int
}

type Command struct {
	P   *Player
	Msg string
}

func (g *Game) Play() {
	q := make(chan *Player, 1000)
	num := 0
	go qwatcher(q, g)
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	check_err(err)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting new connection %v\n", err)
		} else {
			num++
			go handleNewConn(conn, q, num)
		}
	}
}

type Player struct {
	Name   string
	conn   net.Conn
	Num    int
	Active bool
	Id     int
}

func (p *Player) SendMsg(s string) (n int, err error) {
	bs := []byte(s)
	n, err = p.conn.Write(bs)
	return n, err
}

func (p *Player) EndGame(msg string) {
	p.SendMsg(msg)
	p.conn.Close()
}

func (p *Player) RecvMsg() (s string, err error) {
	buff := make([]byte, 100)
	n, err := p.conn.Read(buff)
	if err != nil {
		return "", err
	}

	cutoff := 0 // for handling different line terminator styles
	if n > 0 {
		if buff[n-1] == 10 {
			cutoff = 1 // netcat uses just a LF char (ascii 10)
		}
	}
	if n > 1 {
		if buff[n-2] == 13 {
			cutoff = 2 // telnet uses CR LF (ascii 13 10)
		}
	}
	if n-cutoff > 0 {
		return string(buff[:n-cutoff]), nil
	} else {
		log.Println("trying again")
		return p.RecvMsg()
	}
}

func handleNewConn(conn net.Conn, q chan *Player, num int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprintf("Had an error with %v\n", conn.RemoteAddr()))
		}
	}()
	p := &Player{"", conn, 0, true, num}
	p.SendMsg("Type your name, then hit enter: ")

	name, err := p.RecvMsg()
	if err != nil {
		log.Printf("Could not get player name, err=%v\n", err)
		return
	}
	p.Name = name
	_, err = p.SendMsg(fmt.Sprintf("Thanks %v! We'll get you into the game as soon as possible.\n", p.Name))
	if err != nil {
		log.Printf("Error: %v - while sending to player - could not add '%v' to the queue.\n", err, name)
		return
	}
	q <- p
}

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func handlePlayer(p *Player, coms chan Command, quit chan bool) {
	var err error
	var msg string
	var com Command
	for {
		err = p.conn.SetReadDeadline(time.Now().Add(time.Second))
		msg, err = p.RecvMsg()
		if isTimeout(err) {
			err = nil
			select {
			case <-quit:
				return
			default:
				continue
			}
		}
		if err != nil {
			log.Printf("problem with player %v, error=%v\n", p, err)
			p.EndGame("Connection trouble - kicking you :)\n")
			msg = "quit"
			com = Command{p, msg}
			p.Active = false
			coms <- com
			return
		}
		com = Command{p, msg}
		coms <- com
		select {
		case _, _ = <-quit:
			return
		default:
			continue
		}
	}
}

// TODO configurable number of players per game have a minimum number, a maximum number, and a timeout
// the game doesn't start until the minimum number has been reached
// the game does not start until the timeout has expired, or the maximum number has been reached
// every time a player joins the timeout is (at least partially) reset
// there is a max waiting time for the player who has been waiting the longest (to guard against lots of people joining and leaving)
func qwatcher(q chan *Player, g *Game) {
	pslice := make([]*Player, g.NumPlayers)
	names := make([]string, g.NumPlayers)
	for {
		for i := 0; i < g.NumPlayers; i++ {
			// TODO verify player connection is still alive / player still active
			pslice[i] = <-q
			names[i] = pslice[i].Name
		}

		log.Printf("I got players named %v\n", strings.Join(names, ", "))
		go g.PlayGame(q, pslice)
	}
}

// Wraps Game.Playfunc with error handling so that one misbehaving Playfunc won't crash the whole app
func (g *Game) PlayGame(q chan *Player, ps []*Player) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering in PlayGame: r=%v", r)
			for _, p := range ps {
				p.conn.Close()
			}
		}
	}()
	commands_chan := make(chan Command, 100)
	quit_chan := make(chan bool, 0)
	for _, p := range ps {
		go handlePlayer(p, commands_chan, quit_chan)
	}
	g.Playfunc(ps, commands_chan, quit_chan)
	close(quit_chan)
	for _, p := range ps {
		if p.Active {
			q <- p
		}
	}
}

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}
