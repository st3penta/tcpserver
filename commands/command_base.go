package commands

import "encoding/binary"

func ParseCommand(body []byte) Command {
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

/*
Sample login packet (exluded length):

	01							    0, 1 - Version
		0001						1, 2 - Command
			00000001				3, 4 - CorrelationId
					0003			7, 2 - Username length
						617364	    9, 3 - Username
*/
func parseMetadata(body []byte) Metadata {
	return Metadata{
		version:       body[0],
		cmdCode:       binary.BigEndian.Uint16(body[1:3]),
		correlationId: binary.BigEndian.Uint32(body[3:7]),
	}
}
