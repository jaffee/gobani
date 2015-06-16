package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		panic(err)
	}

	msgs := make(chan string, 200)
	var read_string string
	reader := bufio.NewReader(os.Stdin)
	go readConn(conn, msgs)
	for {
		select {
		case read_string = <-msgs:
			fmt.Println("read")
		default:
			read_string = ""
		}
		if read_string != "" {
			fmt.Println("Received:")
			fmt.Println(read_string)
		}

		if reader.Buffered() > 0 {
			fmt.Println("something buffered")
			b, err := reader.ReadByte()
			if err != nil {
				panic(err)
			}
			conn.Write([]byte{b})
		} else {
			fmt.Println("nothing buffered")
		}
		time.Sleep(time.Second * 1)

	}

}

func readConn(conn net.Conn, msgs chan string) {
	recBuff := make([]byte, 1000)
	for {
		n, err := conn.Read(recBuff)
		fmt.Println("hey")
		if err != nil {
			panic(err)
		}
		msgs <- string(recBuff[:n])
		zeroBuff(recBuff)
	}
}

func zeroBuff(abuff []byte) {
	for i := 0; i < len(abuff); i++ {
		abuff[i] = 0
	}
}
