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

var height = 20
var width = 70
var OPEN = " "
var BLOCKED = "#"
var GOLD = "G"

type position struct {
	x int
	y int
}

type Player struct {
	game.Player
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
	players       map[int]*Player
	gold_position position
}

func goldrace(ps []*game.Player, commands_chan chan game.Command, quit_chan chan bool) {
	players := makePlayers(ps)
	welcome(players)
	b := makeBattlefield(height, width, players)
	for _, p := range players {
		p.SendMsg(b.toString())
	}

	for {
		com := <-commands_chan
		if com.Msg == "quit" {
			delete(players, com.P.Id)
			continue
		}
		com_player := players[com.P.Id]
		com_player.Move(com.Msg, b)
		for _, p := range players {
			p.SendMsg(b.toString())
		}
		if com_player.pos == b.gold_position {
			for _, p := range players {
				if p == com_player {
					p.SendMsg("You Win!\n")
				} else {
					p.SendMsg("You Lose!\n")
				}
			}
			time.Sleep(time.Second * 3)
			return
		}
	}
}

func makePlayers(ps []*game.Player) map[int]*Player {
	players := make(map[int]*Player)
	for _, p := range ps {
		players[p.Id] = &Player{*p, position{}}
	}
	return players
}

func welcome(players map[int]*Player) {
	for i, p := range players {
		p.Num = i
		p.SendMsg(fmt.Sprintf("You are player %v - use 'w', 's', 'a', and 'd' to get to the Gold!\n", i))
		p.SendMsg("'#' characters are walls, 'G' represents the gold. Numbers are other players.\n")
	}
	// Countdown clock
	for count := 5; count >=0; count-- {
		for _, p := range players {
			if count == 0 {
				p.SendMsg("\r")
			} else {
				p.SendMsg(fmt.Sprintf("\rStarting in %v...", count))
			}
		}
		if count > 0 {
			time.Sleep(time.Second * 1)
		}
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

func printPlayers(players map[int]*Player) {
	for i, p := range players {
		fmt.Printf("%v: %v, ", i, *p)
	}
	fmt.Println("")
}

func makeBattlefield(height int, width int, players map[int]*Player) battlefield {
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
		// can really only support up to 10 players because of this
		b.rep[p.pos.y][p.pos.x] = strconv.Itoa(p.Num % 10)
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
