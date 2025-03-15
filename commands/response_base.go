package commands

import (
	"encoding/binary"
	"io"
)

type Response struct {
	responseLength uint32
	Metadata
	statusCode uint16
}

func (r *Response) Write(out io.Writer) {
	binary.Write(out, binary.BigEndian, r.responseLength)
	binary.Write(out, binary.BigEndian, r.Metadata.version)
	binary.Write(out, binary.BigEndian, r.Metadata.cmdCode)
	binary.Write(out, binary.BigEndian, r.Metadata.correlationId)
	binary.Write(out, binary.BigEndian, r.statusCode)
}
