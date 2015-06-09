package main

import (
	"fmt"
	"github.com/jaffee/gobani/game"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	g := game.Game{goldrace, 2}
	g.Play()
}

var height = 15
var width = 20
var OPEN = " "
var BLOCKED = "#"
var GOLD = "G"

type command struct {
	p   *Player
	msg string
}

type position struct {
	x int
	y int
}

type Player struct {
	*game.Player
	pos position
}

func (p *Player) Move(msg string, b battlefield) {
	newpos := position{p.pos.x, p.pos.y}
	if msg == "w" && newpos.y-1 >= 0 {
		newpos.y = newpos.y - 1
	} else if msg == "s" && newpos.y+1 < b.height {
		newpos.y = newpos.y + 1
	} else if msg == "d" && newpos.x+1 < b.width {
		newpos.x = newpos.x + 1
	} else if msg == "a" && newpos.x-1 >= 0 {
		newpos.x = newpos.x - 1
	}

	if b.Rep()[newpos.y][newpos.x] == OPEN || b.Rep()[newpos.y][newpos.x] == GOLD {
		p.pos = newpos
	}
}

type battlefield struct {
	height int
	width  int
	rep    [][]string
	// TODO add obstacles to battlefield
	players       []*Player
	gold_position position
}

func goldrace(ps []game.Player) {
	players := makePlayers(ps)
	welcome(players)
	b := makeBattlefield(height, width, players)
	commands_chan := make(chan command, 100)
	quit_chan := make(chan bool, len(ps)+1)
	for _, p := range players {
		go handlePlayer(p, commands_chan, quit_chan)
		p.SendMsg(b.toString())
	}

	for {
		com := <-commands_chan
		if com.msg == "quit" {
			players = append(players[:com.p.Num], players[com.p.Num+1:]...)
			// TODO handle this everywhere - don't re-enqueue quit players etc.
		}
		com.p.Move(com.msg, b)
		for _, p := range players {
			p.SendMsg(b.toString())
		}
		if com.p.pos == b.gold_position {
			for _, p := range players {
				if p == com.p {
					p.SendMsg("You Win!\n")
				} else {
					p.SendMsg("You Lose!\n")
				}
			}
			for i := 0; i < len(ps); i++ {
				quit_chan <- true
			}
			time.Sleep(time.Second * 3)
			return
		}
	}
}

func makePlayers(ps []game.Player) []*Player {
	players := make([]*Player, len(ps))
	for i := 0; i < len(ps); i++ {
		players[i] = &Player{&ps[i], position{}}
	}
	return players
}

func handlePlayer(p *Player, coms chan command, quit chan bool) {
	for {
		msg, err := p.RecvMsg()
		if err != nil {
			log.Printf("problem with player %v, error=%v\n", p, err)
			p.EndGame("Connection trouble - kicking you :)\n")
			msg = "quit"
			com := command{p, msg}
			coms <- com
			return
		}
		com := command{p, msg}
		coms <- com
		select {
		case <-quit:
			return
		default:
			continue
		}
	}
}

func welcome(players []*Player) {
	for i, p := range players {
		p.Num = i
		p.SendMsg(fmt.Sprintf("You are player %v - use 'w', 's', 'a', and 'd' to get to the Gold!\n", i))
	}
}

func randOpenPos(height, width int, rep [][]string) position {
	numOpen := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if rep[i][j] == OPEN {
				numOpen++
			}
		}
	}
	posNum := rand.Intn(numOpen)
	numOpen = 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if rep[i][j] == OPEN {
				if posNum == numOpen {
					return position{j, i}
				}
				numOpen++
			}
		}
	}
	panic("No open positions")
}

func makeBattlefield(height int, width int, players []*Player) battlefield {
	rep := make([][]string, height)
	for i := 0; i < height; i++ {
		rep[i] = make([]string, width)
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				rep[i][j] = BLOCKED
			} else {
				rep[i][j] = OPEN
			}
		}
	}

	for _, p := range players {
		p.pos = randOpenPos(height, width, rep)
	}

	goldPos := randOpenPos(height, width, rep)
	rep[goldPos.y][goldPos.x] = GOLD

	return battlefield{height, width, rep, players, goldPos}
}

func (b *battlefield) Rep() [][]string {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				b.rep[i][j] = BLOCKED
			} else {
				b.rep[i][j] = OPEN
			}
		}
	}
	for _, p := range b.players {
		b.rep[p.pos.y][p.pos.x] = strconv.Itoa(p.Num)
	}
	b.rep[b.gold_position.y][b.gold_position.x] = GOLD
	return b.rep
}

func (b *battlefield) toString() string {
	// TODO maybe cache whether string rep needs to be updated
	b.rep = b.Rep()
	lines := make([]string, b.height)
	for i, row := range b.rep {
		lines[i] = strings.Join(row, "")
	}
	return strings.Join(lines, "\n") + "\n"
}
