package game

import (
	"fmt"
	"net"
	"strings"
)

type Gameplay func(...Player)

type Game struct {
	Playfunc   Gameplay
	NumPlayers int
}

type Player struct {
	Name string
	conn net.Conn
}

func (g *Game) Play() {
	q := make(chan Player, 1000)
	go qwatcher(q, g)
	l, err := net.Listen("tcp", ":54321")
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

func handleNewConn(conn net.Conn, q chan Player) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("Had an error with %v\n", conn.RemoteAddr()))
		}
	}()
	p := Player{"", conn}
	p.SendMsg("Type your name, then hit enter: ")
	p.Name = p.RecvMsg()
	p.SendMsg(fmt.Sprintf("Thanks %v! We'll get you into the game as soon as possible.\n", p.Name))
	q <- p
}

func qwatcher(q chan Player, g *Game) {
	pslice := make([]Player, g.NumPlayers)
	names := make([]string, g.NumPlayers)
	for {
		for i := 0; i < g.NumPlayers; i++ {
			pslice[i] = <-q
			names[i] = pslice[i].Name
		}

		fmt.Printf("I got players named %v\n", strings.Join(names, ", "))
		go g.PlayGame(q, pslice...)
	}
}

// Wraps Game.Playfunc with error handling so that one misbehaving Playfunc won't crash the whole app
func (g *Game) PlayGame(q chan Player, ps ...Player) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			for _, p := range ps {
				p.conn.Close()
			}
		}
	}()
	g.Playfunc(ps...)
	for _, p := range ps {
		q <- p
	}
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
	if n > 2 {
		return string(buff[:n-2])
	} else {
		return p.RecvMsg()
	}
}

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}
