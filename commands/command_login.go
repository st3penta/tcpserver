package commands

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	LoginCommandCode uint16 = 0x01
)

type LoginCommand struct {
	metadata Metadata
	username string
}

func NewLoginCommand(
	metadata Metadata,
	stream io.Reader,
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
	}

	lc.print()

	return lc, nil
}

func (lc *LoginCommand) Process() (*Response, error) {
	return &Response{
		version:       lc.metadata.version,
		correlationID: lc.metadata.correlationId,
		statusCode:    ResponseStatusCodeOK,
	}, nil
}

func (lc *LoginCommand) print() {
	fmt.Println("-----")
	fmt.Println("Login")
	fmt.Printf("   version: %d\n", lc.metadata.version)
	fmt.Printf("   correlationId: %d\n", lc.metadata.correlationId)
	fmt.Printf("   username: %s\n", lc.username)
	fmt.Println("-----")
}
