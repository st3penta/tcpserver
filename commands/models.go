package commands

import (
	"encoding/binary"
	"io"
)

// Command
type Metadata struct {
	version       byte
	cmdCode       uint16
	correlationId uint32
}

type Command interface {
	Process(io.Writer)
	Print()
}

// Response
type Response interface {
	Marshal() []byte
}

type BaseResponse struct {
	Metadata
	StatusCode uint16
}

func (r *BaseResponse) Write(out io.Writer) []byte {
	respBytes := make([]byte, 13)

	binary.Write(out, binary.BigEndian, uint32(9))
	binary.Write(out, binary.BigEndian, r.version)
	binary.Write(out, binary.BigEndian, r.Metadata.cmdCode)
	binary.Write(out, binary.BigEndian, r.Metadata.correlationId)
	binary.Write(out, binary.BigEndian, r.StatusCode)

	return respBytes
}
