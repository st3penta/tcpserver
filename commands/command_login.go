package commands

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	LoginResponseLength = 9
	LoginCommandCode    = 3
)

var (
	ErrLoginCommandTooShort  = errors.New("malformed login command: message is too short")
	ErrInvalidUsernameLength = errors.New("malformed login command: invalid username length")
)

type LoginCommand struct {
	metadata Metadata
	username string
}

// Sample login packet (exluded length and metadata):
// ..............0003			7, 2 - Username length
// ..................617364		9, 3 - Username
func NewLoginCommand(
	body []byte,
) (*LoginCommand, error) {

	metadata, mErr := parseMetadata(body)
	if mErr != nil {
		return nil, mErr
	}

	if len(body) < 9 {
		return nil, ErrLoginCommandTooShort
	}

	usernameLen := binary.BigEndian.Uint16(body[7:9])

	if len(body) != int(9+usernameLen) {
		return nil, ErrInvalidUsernameLength
	}

	lc := &LoginCommand{
		metadata: metadata,
		username: string(body[9:(9 + usernameLen)]),
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
