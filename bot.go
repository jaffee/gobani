package main

import (
	"fmt"
	"net"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:54321")
	checkErr(err)

	conn.Write([]byte("bot\nrock\n"))
	recBuff := make([]byte, 190)
	conn.Read(recBuff)
	fmt.Println(string(recBuff))
	conn.Write([]byte("rock\n"))
	fmt.Println(string(recBuff))

}
