package commands

import (
	"encoding/binary"
	"fmt"
	"io"
)

type LoginCommand struct {
	Metadata
	CorrelationId uint32
	UsernameLen   uint16
	Username      string
}

type LoginResponse = BaseResponse

func NewLoginCommand(
	body []byte,
) *LoginCommand {

	version := body[0]
	cmdCode := binary.BigEndian.Uint16(body[1:3])

	cmdBytes := body[3:]
	correlationId := binary.BigEndian.Uint32(cmdBytes[:4])
	usernameLen := binary.BigEndian.Uint16(cmdBytes[4:6])
	username := string(cmdBytes[6:(6 + usernameLen)])

	lc := &LoginCommand{
		Metadata: Metadata{
			version: version,
			cmdCode: cmdCode,
		},
		CorrelationId: correlationId,
		UsernameLen:   usernameLen,
		Username:      username,
	}

	lc.Print()

	return lc
}

func (lc *LoginCommand) Process(out io.Writer) {
	resp := &LoginResponse{
		Metadata: Metadata{
			version: lc.version,
			cmdCode: 3,
		},
		CorrelationId: lc.CorrelationId,
		StatusCode:    1,
	}

	resp.Write(out)
}

func (lc *LoginCommand) Print() {
	fmt.Println("-----")
	fmt.Println("Login")
	fmt.Println(fmt.Sprintf("version: %d", lc.version))
	fmt.Println(fmt.Sprintf("correlationId: %d", lc.CorrelationId))
	fmt.Println(fmt.Sprintf("usernameLen: %d", lc.UsernameLen))
	fmt.Println(fmt.Sprintf("username: %s", lc.Username))
	fmt.Println("-----")
}
