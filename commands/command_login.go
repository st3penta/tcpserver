package commands

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	LoginResponseLength = uint32(9)
	LoginCommandCode    = uint16(1)

	ErrLoginCommandTooShort  = errors.New("malformed login command: message is too short")
	ErrInvalidUsernameLength = errors.New("malformed login command: invalid username length")
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
		responseLength: LoginResponseLength,
		Metadata: Metadata{
			version:       lc.metadata.version,
			cmdCode:       LoginCommandCode,
			correlationId: lc.metadata.correlationId,
		},
		statusCode: 1,
	}, nil
}

func (lc *LoginCommand) print() {
	fmt.Println("-----")
	fmt.Println("Login")
	fmt.Println(fmt.Sprintf("   version: %d", lc.metadata.version))
	fmt.Println(fmt.Sprintf("   correlationId: %d", lc.metadata.correlationId))
	fmt.Println(fmt.Sprintf("   username: %s", lc.username))
	fmt.Println("-----")
}
