package commands

import (
	"encoding/binary"
	"fmt"
	"io"
)

type LoginCommand struct {
	metadata Metadata
	username string
}

type LoginResponse = BaseResponse

func NewLoginCommand(
	body []byte,
) *LoginCommand {

	/*
		Sample login packet (exluded length):

			01							    0, 1 - Version
				0001						1, 2 - Command
					00000001				3, 4 - CorrelationId
							0003			7, 2 - Username length
								617364	    9, 3 - Username
	*/
	usernameLen := binary.BigEndian.Uint16(body[7:9])

	lc := &LoginCommand{
		metadata: parseMetadata(body),
		username: string(body[9:(9 + usernameLen)]),
	}

	lc.Print()

	return lc
}

func (lc *LoginCommand) Process(out io.Writer) {
	resp := &LoginResponse{
		Metadata: Metadata{
			version:       lc.metadata.version,
			cmdCode:       3,
			correlationId: lc.metadata.correlationId,
		},
		StatusCode: 1,
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
