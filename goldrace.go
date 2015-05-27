package main

import (
	"github.com/jaffee/gobani/game"
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
}

type battlefield struct {
	height int
	width  int
	rep    [][]string
}

func makeBattlefield(height int, width int) battlefield {
	rep := make([][]string, height)
	for i := 0; i < height; i++ {
		rep[i] = make([]string, width)
		for j := 0; j < width; j++ {
			rep[i][j] = "O"
		}
	}

	return battlefield{height, width, rep}
}

func (b *battlefield) toString() string {
	lines := make([]string, b.height)
	for i, row := range b.rep {
		lines[i] = strings.Join(row, "")
	}
	return strings.Join(lines, "\n") + "\n"
}
