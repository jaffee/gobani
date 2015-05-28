package main

import (
	"github.com/jaffee/gobani/game"
	"math/rand"
	"strings"
)

func main() {
	g := game.Game{goldrace, 1}
	g.Play()
}

var height = 10
var width = 10

func goldrace(ps ...game.Player) {
	battlefield := makeBattlefield(height, width)
	ps[0].SendMsg(battlefield.toString())
	for {
		msg := ps[0].RecvMsg()
		if msg == "u" {
			battlefield.Move("P", "u")
		}
		ps[0].SendMsg(battlefield.toString())
	}
}

type battlefield struct {
	height int
	width  int
	rep    [][]string
	px     int
	py     int
	gx     int
	gy     int
}

func makeBattlefield(height int, width int) battlefield {
	rep := make([][]string, height)
	for i := 0; i < height; i++ {
		rep[i] = make([]string, width)
		for j := 0; j < width; j++ {
			rep[i][j] = "O"
		}
	}
	px := rand.Intn(width)
	py := rand.Intn(height)
	gx := rand.Intn(width)
	gy := rand.Intn(height)

	rep[py][px] = "P"
	rep[gy][gx] = "G"

	return battlefield{height, width, rep, px, py, gx, gy}
}

func (b *battlefield) toString() string {
	lines := make([]string, b.height)
	for i, row := range b.rep {
		lines[i] = strings.Join(row, "")
	}
	return strings.Join(lines, "\n") + "\n"
}

func (b *battlefield) Move(thing string, move string) {
	if move == "u" && b.py-1 >= 0 {
		b.rep[b.py][b.px] = "O"
		b.py = b.py - 1
		b.rep[b.py][b.px] = "P"
	}
}
