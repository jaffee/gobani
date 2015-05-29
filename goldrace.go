package main

import (
	"github.com/jaffee/gobani/game"
	"math/rand"
	"strconv"
	"strings"
)

func main() {
	g := game.Game{goldrace, 2}
	g.Play()
}

var height = 10
var width = 10
var OPEN = "O"
var GOLD = "G"

func goldrace(ps ...game.Player) {
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
			ps[winner].SendMsg("You Win!")
		}
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
			rep[i][j] = OPEN
		}
	}
	positions := make([]position, numPlayers)
	for i := 0; i < numPlayers; i++ {
		pos := randPos(height, width)
		positions[i] = pos
		rep[pos.y][pos.x] = strconv.Itoa(i)
	}
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
	if move == "u" && pos.y-1 >= 0 {
		b.rep[pos.y][pos.x] = OPEN
		pos.y = pos.y - 1
	} else if move == "d" && pos.y+1 < b.height {
		b.rep[pos.y][pos.x] = OPEN
		pos.y = pos.y + 1
	} else if move == "r" && pos.x+1 < b.width {
		b.rep[pos.y][pos.x] = OPEN
		pos.x = pos.x + 1
	} else if move == "l" && pos.x-1 >= 0 {
		b.rep[pos.y][pos.x] = OPEN
		pos.x = pos.x - 1
	}
	b.rep[pos.y][pos.x] = strconv.Itoa(pnum)
	b.positions[pnum] = pos
}
