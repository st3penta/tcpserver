package commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	LoginCommandCode uint16 = 0x01
)

type LoginCommand struct {
	metadata Metadata
	username string
	conn     net.Conn
}

func NewLoginCommand(
	metadata Metadata,
	stream io.Reader,
	conn net.Conn,
) (*LoginCommand, error) {

	var usernameLen uint16
	lErr := binary.Read(stream, binary.BigEndian, &usernameLen)
	if lErr != nil {
		return nil, lErr
	}

	username := make([]byte, usernameLen)
	_, uErr := io.ReadFull(stream, username)
	if uErr != nil {
		return nil, uErr
	}

	lc := &LoginCommand{
		metadata: metadata,
		username: string(username),
		conn:     conn,
	}

	lc.print()

	return lc, nil
}

func (lc *LoginCommand) Process(state State) (*Response, error) {

	err := state.Login(lc.conn, lc.username)
	if err != nil {
		return nil, err
	}

	return &Response{
		version:       lc.metadata.version,
		correlationID: lc.metadata.correlationId,
		statusCode:    ResponseStatusCodeOK,
	}, nil
}

func (lc *LoginCommand) print() {
	fmt.Println("-----")
	fmt.Println("Login")
	fmt.Printf("\tversion: %d\n", lc.metadata.version)
	fmt.Printf("\tcorrelationId: %d\n", lc.metadata.correlationId)
	fmt.Printf("\tusername: %s\n", lc.username)
	fmt.Println("-----")
}
