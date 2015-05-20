package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

type player struct {
	Name string
	conn net.Conn
}

var legalMoves = []string{"rock", "paper", "scissors"}

func (p *player) SendMsg(s string) (n int, err error) {
	bs := []byte(s)
	n, err = p.conn.Write(bs)
	return n, err
}

func (p *player) RecvMsg() (s string) {
	buff := make([]byte, 100)
	n, err := p.conn.Read(buff)
	if err != nil {
		fmt.Println("Problem receiving message: ", err)
	}
	if n < 2 {
		fmt.Println("Problem!")
		fmt.Println(string(buff))
	}
	return string(buff[:n-2])
}

type battleground struct {
	ground []int
}

func NewBattleground() {
	bgsize := 100
	bgmaxheight := 100
	ground := make([]int, bgsize)
	for i := 0; i < bgsize; i++ {
		ground[i] = rand.Intn(bgmaxheight)
	}

}

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	q := make(chan player, 100)
	go qwatcher(q)
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

func handleNewConn(conn net.Conn, q chan player) {
	p := player{"", conn}
	p.Name = p.RecvMsg()
	q <- p
}

func qwatcher(q chan player) {
	var p1, p2 player
	for {
		p1 = <-q
		p2 = <-q

		fmt.Printf("I got two players named %v and %v - let's do this shiznizzle!\n", p1.Name, p2.Name)
		theThunderdome(p1, p2)
	}
}

func theThunderdome(p1 player, p2 player) {
	welcome(p1, p2)
	p1score := 0
	p2score := 0
	var p1msg string
	var p2msg string
	for {
		ch1 := make(chan string, 1)
		go getMove(p1, ch1)
		ch2 := make(chan string, 2)
		go getMove(p2, ch2)
		move1 := <-ch1
		move2 := <-ch2
		winner := getWinner(move1, move2)
		if winner == 1 {
			p1score++
			p1msg = "win"
			p2msg = "lose"
		} else if winner == 2 {
			p2score++
			p1msg = "lose"
			p2msg = "win"
		} else {
			p1msg = "tie"
			p2msg = "tie"
		}

		fmt.Printf("Player 1 played %v, and p2 played %v\n", move1, move2)
		p1.SendMsg(fmt.Sprintf("%v:%v played %v and you:%v played %v. You %v.\n", p2.Name, p2score, move2, p1score, move1, p1msg))
		p2.SendMsg(fmt.Sprintf("%v:%v played %v and you:%v played %v. You %v.\n", p1.Name, p1score, move1, p2score, move2, p2msg))
	}
}

func getWinner(m1 string, m2 string) int {
	m1index := stringIndex(m1, legalMoves)
	m2index := stringIndex(m2, legalMoves)
	if m1index == m2index {
		return 0
	} else if (m1index+1)%len(legalMoves) == m2index {
		return 2
	} else {
		return 1
	}

}

func stringIndex(s string, sslice []string) int {
	for i, str := range sslice {
		if s == str {
			return i
		}
	}
	return -1
}

func getMove(p player, ch chan string) {
	for {
		ans := p.RecvMsg()
		for i := 0; i < len(legalMoves); i++ {
			if ans == legalMoves[i] {
				ch <- ans
				return
			}
		}
		p.SendMsg(fmt.Sprintf("That's not a legal move! Try one of these: %v\n", strings.Join(legalMoves, ", ")))
	}
}

func welcome(p1 player, p2 player) {
	welcome_msg := "%v! You are now entering... the THUNDERDOME. Your opponent, %v, is going to %v.\n"
	p1.SendMsg(fmt.Sprintf(welcome_msg, p1.Name, p2.Name, randomTaunt()))
	p2.SendMsg(fmt.Sprintf(welcome_msg, p2.Name, p1.Name, randomTaunt()))

}

func randomTaunt() string {
	taunts := []string{"eat your socks", "rub honey all over you",
		"strangle you with your umbilical cord", "curbstomp your teddy bear",
		"hit you, like, really hard", "bodyslam your grandmother", "forcefeed you kale"}
	i := rand.Intn(len(taunts))
	return taunts[i]
}
