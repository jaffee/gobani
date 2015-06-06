package main

import (
	"fmt"
	"github.com/jaffee/gobani/game"
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

func (p *Player) Move(msg string, height int, width int) {
	if msg == "w" && p.pos.y-1 >= 0 {
		p.pos.y = p.pos.y - 1
	} else if msg == "s" && p.pos.y+1 < height {
		p.pos.y = p.pos.y + 1
	} else if msg == "d" && p.pos.x+1 < width {
		p.pos.x = p.pos.x + 1
	} else if msg == "a" && p.pos.x-1 >= 0 {
		p.pos.x = p.pos.x - 1
	}
}

type battlefield struct {
	height        int
	width         int
	rep           [][]string
	players       []*Player
	gold_position position
	// cache whether rep has been updated
}

func goldrace(ps ...game.Player) {
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
		com.p.Move(com.msg, b.height, b.width)
		for _, p := range players {
			p.SendMsg(b.toString())
		}
		fmt.Printf("%v, %v\n", com.p.pos, b.gold_position)
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

// TODO realtime instead of turn based
func handlePlayer(p *Player, coms chan command, quit chan bool) {
	for {
		msg := p.RecvMsg()
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

func randPos(height int, width int) (pos position) {
	pos = position{rand.Intn(width), rand.Intn(height)}
	return
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
		p.pos = randPos(height, width)
	}
	// TODO detect collisions

	goldPos := randPos(height, width)
	rep[goldPos.y][goldPos.x] = GOLD

	return battlefield{height, width, rep, players, goldPos}
}

func (b *battlefield) toString() string {
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

	lines := make([]string, b.height)
	for i, row := range b.rep {
		lines[i] = strings.Join(row, "")
	}
	return strings.Join(lines, "\n") + "\n"
}

// func (b *battlefield) Move(com command) {
// 	pos := b.positions[com.p.Num]
// 	move := com.msg
// 	if move == "w" && pos.y-1 >= 0 {
// 		b.rep[pos.y][pos.x] = OPEN
// 		pos.y = pos.y - 1
// 	} else if move == "s" && pos.y+1 < b.height {
// 		b.rep[pos.y][pos.x] = OPEN
// 		pos.y = pos.y + 1
// 	} else if move == "d" && pos.x+1 < b.width {
// 		b.rep[pos.y][pos.x] = OPEN
// 		pos.x = pos.x + 1
// 	} else if move == "a" && pos.x-1 >= 0 {
// 		b.rep[pos.y][pos.x] = OPEN
// 		pos.x = pos.x - 1
// 	}
// 	b.rep[pos.y][pos.x] = strconv.Itoa(pnum)
// 	b.positions[com.p.Num] = pos
// }
