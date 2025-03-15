package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"tcpserver/commands"
)

type Server struct {
	port        int
	loggedUsers map[string]bool
}

func NewServer(port int) *Server {
	return &Server{
		port:        port,
		loggedUsers: map[string]bool{},
	}
}

func (s *Server) Start() {

	// Start listening on the specified port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Accept incoming connections (in async, so that we can accept multiple connections)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {

	for {
		fmt.Println("Waiting for new messages...")

		msgLen, lenErr := s.readCmdLength(conn)
		if lenErr != nil {
			break
		}

		body, bodyErr := s.readCmdBody(conn, msgLen)
		if bodyErr != nil {
			break
		}

		cmd := commands.ParseCommand(body)

		cmd.Process(conn)
	}
}

func (s *Server) readCmdLength(conn net.Conn) (uint32, error) {
	msgLenBytes := make([]byte, 4)
	_, err := conn.Read(msgLenBytes)
	if err != nil {
		fmt.Println("Error while reading msg length: ", err)
		return 0, err
	}
	return binary.BigEndian.Uint32(msgLenBytes), nil
}

func (s *Server) readCmdBody(conn net.Conn, msgLen uint32) ([]byte, error) {
	body := make([]byte, msgLen)
	_, err := conn.Read(body)
	if err != nil {
		fmt.Println("Error while reading connection", err)
		return nil, err
	}
	fmt.Println(fmt.Sprintf("Body: %X", body))
	return body, nil
}
