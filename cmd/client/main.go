package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	addr = flag.String("addr", "localhost:9876", "server address")
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		line, err := stdinReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		if _, err := conn.Write([]byte(line)); err != nil {
			panic(err)
		}

		response, err := connReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		fmt.Print(response)
	}
}
