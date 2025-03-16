package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
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

func ParseCommand(stream io.Reader) (Command, error) {

	body, bErr := readFieldWithLength(stream)
	if bErr != nil {
		return nil, bErr
	}

	bodyStream := bytes.NewBuffer(body)

	metadata, mErr := parseMetadata(bodyStream)
	if mErr != nil {
		return nil, mErr
	}

	switch metadata.cmdCode {
	case LoginCommandCode:
		return NewLoginCommand(*metadata, bodyStream)
	case CorrelationIDTestCommandCode:
		return NewCorrelationIDTestCommand(*metadata, bodyStream)
	// case MessageCommandCode:
	// 	return s.parseCmdMessage(version, cmdCode, body[3:])
	default:
		return nil, ErrUnknownCommand
	}
}

func readFieldWithLength(stream io.Reader) ([]byte, error) {
	var fieldLen uint32
	err := binary.Read(stream, binary.BigEndian, &fieldLen)
	if err != nil {
		return nil, err
	}

	field := make([]byte, fieldLen)
	_, err = io.ReadFull(stream, field)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func parseMetadata(stream io.Reader) (*Metadata, error) {

	var version byte
	vErr := binary.Read(stream, binary.BigEndian, &version)
	if vErr != nil {
		return nil, vErr
	}

	var cmdCode uint16
	cErr := binary.Read(stream, binary.BigEndian, &cmdCode)
	if cErr != nil {
		return nil, cErr
	}

	var correlationId uint32
	crErr := binary.Read(stream, binary.BigEndian, &correlationId)
	if crErr != nil {
		return nil, crErr
	}

	return &Metadata{
		version:       version,
		cmdCode:       cmdCode,
		correlationId: correlationId,
	}, nil
}
