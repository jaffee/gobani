package main

import (
	"fmt"
	"net"
)

type player struct {
	name string
	conn net.Conn
}

type battleground struct {
	ground []int
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
	buff := make([]byte, 100)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Printf("Error reading from conn remote addr:%v, err:%v\n", conn.RemoteAddr(), err)
	}
	p := player{string(buff[:n-2]), conn}
	q <- p
}

func qwatcher(q chan player) {
	var p1, p2 player
	for {
		p1 = <-q
		p2 = <-q

		fmt.Printf("I got two players named %v and %v - let's do this shiznizzle!\n", p1.name, p2.name)

	}
}
