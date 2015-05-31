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

func goldrace(ps ...game.Player) {
	welcome(ps)
	b := makeBattlefield(height, width, len(ps))
	for _, p := range ps {
		p.SendMsg(b.toString())
	}
	for {
		winner := -1
		for i, player := range ps {
			msg := player.RecvMsg()
			b.Move(i, msg)
			pos := b.positions[i]
			if pos == b.gold_position {
				winner = i
			}
		}
		for _, player := range ps {
			player.SendMsg(b.toString())
		}
		if winner != -1 {
			for i, p := range ps {
				if i == winner {
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

// func handlePlayer(p game.Player, msgs chan string, quit chan bool) {
// 	for {
// 		msg := p.RecvMsg()
// 		msgs <- msg
// 		select {
// 		case doquit := <-quit:
// 			return
// 		default:
// 			continue
// 		}
// 	}
// }

func welcome(ps []game.Player) {
	for i, p := range ps {
		p.SendMsg(fmt.Sprintf("You are player %v - use 'w', 's', 'a', and 'd' to get to the Gold!\n", i))
	}
}

type battlefield struct {
	height        int
	width         int
	rep           [][]string
	positions     []position
	gold_position position
}

type position struct {
	x int
	y int
}

func randPos(height int, width int) (ret position) {
	ret = position{rand.Intn(width), rand.Intn(height)}
	return
}

func makeBattlefield(height int, width int, numPlayers int) battlefield {
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
	// TODO maze generator

	positions := make([]position, numPlayers)
	for i := 0; i < numPlayers; i++ {
		pos := randPos(height, width)
		positions[i] = pos
		rep[pos.y][pos.x] = strconv.Itoa(i)
	}
	// TODO detect collisions

	goldPos := randPos(height, width)
	rep[goldPos.y][goldPos.x] = GOLD

	return battlefield{height, width, rep, positions, goldPos}
}

func (b *battlefield) toString() string {
	lines := make([]string, b.height)
	for i, row := range b.rep {
		lines[i] = strings.Join(row, "")
	}
	return strings.Join(lines, "\n") + "\n"
}

func (b *battlefield) Move(pnum int, move string) {
	pos := b.positions[pnum]
	if move == "w" && pos.y-1 >= 0 {
		b.rep[pos.y][pos.x] = OPEN
		pos.y = pos.y - 1
	} else if move == "s" && pos.y+1 < b.height {
		b.rep[pos.y][pos.x] = OPEN
		pos.y = pos.y + 1
	} else if move == "d" && pos.x+1 < b.width {
		b.rep[pos.y][pos.x] = OPEN
		pos.x = pos.x + 1
	} else if move == "a" && pos.x-1 >= 0 {
		b.rep[pos.y][pos.x] = OPEN
		pos.x = pos.x - 1
	}
	b.rep[pos.y][pos.x] = strconv.Itoa(pnum)
	b.positions[pnum] = pos
}
