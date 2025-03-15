package commands

import (
	"encoding/binary"
	"errors"
)

var (
	ErrMalformedMetadata = errors.New("malformed metadata")
	ErrMalformedCommand  = errors.New("malformed command")
	ErrUnknownCommand    = errors.New("unknown command")
)

type Command interface {
	Process() (*Response, error)
}

type Metadata struct {
	version       byte
	cmdCode       uint16
	correlationId uint32
}

func ParseCommand(body []byte) (Command, error) {
	if len(body) < 3 {
		return nil, ErrMalformedCommand
	}

	cmdCode := binary.BigEndian.Uint16(body[1:3])

	switch cmdCode {
	case 1:
		return NewLoginCommand(body)
	// case 2:
	// 	return s.parseCmdMessage(version, cmdCode, body[3:])
	default:
		return nil, ErrUnknownCommand
	}
}

// Sample metadata (offset, length - description):
// 01............   0, 1 - Version
// ..0001........	1, 2 - Command
// ......00000001	3, 4 - CorrelationId
func parseMetadata(body []byte) (Metadata, error) {
	if len(body) < 7 {
		return Metadata{}, ErrMalformedMetadata
	}

	return Metadata{
		version:       body[0],
		cmdCode:       binary.BigEndian.Uint16(body[1:3]),
		correlationId: binary.BigEndian.Uint32(body[3:7]),
	}, nil
}
