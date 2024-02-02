package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/asynched/gokvdb/internal/database"
	"github.com/asynched/gokvdb/internal/server"
)

var (
	addr = flag.String("addr", "localhost:9876", "server address")
)

func init() {
	log.SetFlags(0)

	log.Println(` _________________
|# :           : #|
|  :           :  |
|  :           :  |
|  :           :  |	GO KVDB
|  :___________:  |	Version 0.0.1
|     _________   |	https://github.com/asynched/gokvdb
|    | __      |  |
|    ||  |     |  |
\____||__|_____|__|`)
	log.Println("")

	log.SetFlags(log.Ltime | log.Ldate | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[%d] [gokvdb] ", os.Getpid()))
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatalf("message='Failed to start server' error='%s'\n", err)
	}

	database := database.New()
	server := server.New(listener, database)

	log.Println("message='server started' port=9876")
	if err := server.Run(); err != nil {
		log.Fatalf("message='server failed' error='%s'\n", err)
	}
}
