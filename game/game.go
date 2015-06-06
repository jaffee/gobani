package game

import (
	"fmt"
	"net"
	"strings"
)

type Gameplay func([]Player)

type Game struct {
	Playfunc   Gameplay
	NumPlayers int
}

func (g *Game) Play() {
	q := make(chan Player, 1000)
	go qwatcher(q, g)
	l, err := net.Listen("tcp", ":80")
	check_err(err)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting new connection %v\n", err)
		} else {
			go handleNewConn(conn, q)
		}
	}
}

type Player struct {
	Name string
	conn net.Conn
	Num  int
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

func (p *Player) RecvMsg() (s string) {
	buff := make([]byte, 100)
	n, err := p.conn.Read(buff)
	if err != nil {
		fmt.Println("Problem receiving message: ", err.Error())
		panic(err)
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
	if n > 0 {
		return string(buff[:n-cutoff])
	} else {
		fmt.Println("trying again")
		return p.RecvMsg()
	}
}

func handleNewConn(conn net.Conn, q chan Player) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("Had an error with %v\n", conn.RemoteAddr()))
		}
	}()
	p := Player{"", conn, 0}
	p.SendMsg("Type your name, then hit enter: ")
	p.Name = p.RecvMsg()
	p.SendMsg(fmt.Sprintf("Thanks %v! We'll get you into the game as soon as possible.\n", p.Name))
	q <- p
}

// TODO configurable number of players per game have a minimum number, a maximum number, and a timeout
// the game doesn't start until the minimum number has been reached
// the game does not start until the timeout has expired, or the maximum number has been reached
// every time a player joins the timeout is (at least partially) reset
// there is a max waiting time for the player who has been waiting the longest (to guard against lots of people joining and leaving)
func qwatcher(q chan Player, g *Game) {
	pslice := make([]Player, g.NumPlayers)
	names := make([]string, g.NumPlayers)
	for {
		for i := 0; i < g.NumPlayers; i++ {
			pslice[i] = <-q
			names[i] = pslice[i].Name
		}

		fmt.Printf("I got players named %v\n", strings.Join(names, ", "))
		go g.PlayGame(q, pslice)
	}
}

// Wraps Game.Playfunc with error handling so that one misbehaving Playfunc won't crash the whole app
func (g *Game) PlayGame(q chan Player, ps []Player) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			for _, p := range ps {
				p.conn.Close()
			}
		}
	}()
	g.Playfunc(ps)
	for _, p := range ps {
		q <- p
	}
}

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}
