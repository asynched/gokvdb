package server

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/asynched/gokvdb/internal/commands"
	"github.com/asynched/gokvdb/internal/database"
)

type Server struct {
	listener net.Listener
	database *database.Database
}

func New(listener net.Listener, database *database.Database) *Server {
	return &Server{
		listener: listener,
		database: database,
	}
}

func (server *Server) Run() error {
	for {
		conn, err := server.listener.Accept()

		if err != nil {
			continue
		}

		log.Printf("message='Connection established' address='%s'\n", conn.RemoteAddr())
		go server.handle(conn)
	}
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()
	defer log.Printf("message='Connection closed' address='%s'", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		command, err := reader.ReadString('\n')

		if err != nil {
			return
		}

		command = strings.TrimSpace(command)

		cmd, err := commands.Parse(command)

		if err != nil {
			if _, err := writer.WriteString(err.Error() + "\n"); err != nil {
				return
			}

			if err := writer.Flush(); err != nil {
				return
			}

			continue
		}

		result, err := cmd.Apply(server.database)

		if err != nil {
			if _, err := writer.WriteString(err.Error() + "\n"); err != nil {
				return
			}

			if err := writer.Flush(); err != nil {
				return
			}

			continue
		}

		if _, err := writer.WriteString(result + "\n"); err != nil {
			return
		}

		if err := writer.Flush(); err != nil {
			return
		}
	}
}
