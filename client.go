package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net"
	"os"
	"time"
)

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	stdinMsgs := make(chan byte)
	go readBuff(reader, stdinMsgs)

	msgs := make(chan string, 200)
	var read_string string
	go readConn(conn, msgs)

	for {
		select {
		case read_string = <-msgs:
			// fmt.Println("read")
		default:
			read_string = ""
		}
		if read_string != "" {
			// fmt.Println("Received:")
			fmt.Println(read_string)
		}

		select {
		case read_byte := <-stdinMsgs:
			_, err := conn.Write([]byte{read_byte})
			if err != nil {
				panic(err)
			}
			// fmt.Println("wrote to conn", n)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

}

func readConn(conn net.Conn, msgs chan string) {
	recBuff := make([]byte, 1000)
	for {
		n, err := conn.Read(recBuff)
		if err != nil {
			panic(err)
		}
		// fmt.Println("Read bytes from conn ", n)
		msgs <- string(recBuff[:n])
		zeroBuff(recBuff)
	}
}

func readBuff(reader *bufio.Reader, msgs chan byte) {
	for {
		b, err := reader.ReadByte()
		// fmt.Println("Read byte ", b)
		if err != nil {
			panic(err)
		}
		msgs <- b
		select {
		case <-msgs:
			break
		default:
			continue
		}
	}
}

func zeroBuff(abuff []byte) {
	for i := 0; i < len(abuff); i++ {
		abuff[i] = 0
	}
}
