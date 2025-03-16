package commands

import (
	"encoding/binary"
	"io"
)

const (
	ResponseMsgCode uint16 = 0x03
	ResponseLength  uint32 = 0x0009
)

type Response struct {
	version       byte
	correlationID uint32
	statusCode    uint16
}

func (r *Response) Write(out io.Writer) {
	binary.Write(out, binary.BigEndian, ResponseLength)
	binary.Write(out, binary.BigEndian, r.version)
	binary.Write(out, binary.BigEndian, ResponseMsgCode)
	binary.Write(out, binary.BigEndian, r.correlationID)
	binary.Write(out, binary.BigEndian, r.statusCode)
}
