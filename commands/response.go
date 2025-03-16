package commands

import (
	"encoding/binary"
	"io"
)

const (
	ResponseMsgCode uint16 = 0x03
	ResponseLength  uint32 = 0x0009

	ResponseStatusCodeOK                uint16 = 0x01
	ResponseStatusCodeUserNotFound      uint16 = 0x03
	ResponseStatusCodeUserAlreadyLogged uint16 = 0x04
)

type Response struct {
	version       byte
	correlationID uint32
	statusCode    uint16
}

func (r *Response) Write(out io.Writer) error {
	err := binary.Write(out, binary.BigEndian, ResponseLength)
	if err != nil {
		return err
	}

	err = binary.Write(out, binary.BigEndian, r.version)
	if err != nil {
		return err
	}

	err = binary.Write(out, binary.BigEndian, ResponseMsgCode)
	if err != nil {
		return err
	}

	err = binary.Write(out, binary.BigEndian, r.correlationID)
	if err != nil {
		return err
	}

	return binary.Write(out, binary.BigEndian, r.statusCode)
}
