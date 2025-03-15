package commands

import "encoding/binary"

/*
Sample login packet (exluded length):

	01							0, 1 - Version
		0001						1, 2 - Command
			00000001				3, 4 - CorrelationId
					0003			7, 2 - Username length
						617364	9, 3 - Username
*/
func ParseCmd(body []byte) Command {
	cmdCode := binary.BigEndian.Uint16(body[1:3])

	switch cmdCode {
	case 1:
		return NewLoginCommand(body)
	// case 2:
	// 	return s.parseCmdMessage(version, cmdCode, body[3:])
	default:
		panic("Command not recognized")
	}
}
