package commands

import (
	"encoding/binary"
	"io"
)

const (
	ResponseCode = uint16(3)
)

type Response struct {
	responseLength uint32
	Metadata
	statusCode uint16
}

func (r *Response) Write(out io.Writer) {
	binary.Write(out, binary.BigEndian, r.responseLength)
	binary.Write(out, binary.BigEndian, r.Metadata.version)
	binary.Write(out, binary.BigEndian, ResponseCode)
	binary.Write(out, binary.BigEndian, r.Metadata.correlationId)
	binary.Write(out, binary.BigEndian, r.statusCode)
}
