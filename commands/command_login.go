package commands

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	metadata    Metadata
	responseLen uint32
	username    string
}

type LoginResponse = BaseResponse

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

	/*
		Sample login packet (exluded length):

			01							    0, 1 - Version
				0001						1, 2 - Command
					00000001				3, 4 - CorrelationId
							0003			7, 2 - Username length
								617364	    9, 3 - Username
	*/
	usernameLen := binary.BigEndian.Uint16(body[7:9])

	fmt.Println("body: ", len(body))
	fmt.Println("usernameLen: ", usernameLen)
	if len(body) != int(9+usernameLen) {
		return nil, ErrInvalidUsernameLength
	}

	lc := &LoginCommand{
		metadata:    metadata,
		responseLen: 9,
		username:    string(body[9:(9 + usernameLen)]),
	}

	lc.Print()

	return lc, nil
}

func (lc *LoginCommand) Process(out io.Writer) {
	resp := &LoginResponse{
		responseLength: LoginResponseLength,
		Metadata: Metadata{
			version:       lc.metadata.version,
			cmdCode:       LoginCommandCode,
			correlationId: lc.metadata.correlationId,
		},
		statusCode: 1,
	}

	resp.Write(out)
}

func (lc *LoginCommand) Print() {
	fmt.Println("-----")
	fmt.Println("Login")
	fmt.Println(fmt.Sprintf("version: %d", lc.metadata.version))
	fmt.Println(fmt.Sprintf("correlationId: %d", lc.metadata.correlationId))
	fmt.Println(fmt.Sprintf("username: %s", lc.username))
	fmt.Println("-----")
}
