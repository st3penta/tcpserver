package commands

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	ErrMalformedMetadata     = errors.New("malformed metadata")
	ErrMalformedCommand      = errors.New("malformed command")
	ErrUnknownCommand        = errors.New("unknown command")
	ErrUnsupportedLengthSize = errors.New("unsupported length size")
)

type State interface {
	Login(conn net.Conn, username string) error
	EnqueueMessage(from string, to string, timestamp time.Time, message string) error
}

type Command interface {
	Process(state State) (*Response, error)
}

type Metadata struct {
	version       byte
	cmdCode       uint16
	correlationId uint32
}

func ParseCommand(stream net.Conn) (Command, error) {

	var len uint32
	body, bErr := readFieldWithLength(stream, len)
	if bErr != nil {
		return nil, bErr
	}
	fmt.Printf("### Received command: %X\n", body)

	bodyStream := bytes.NewBuffer(body)

	metadata, mErr := parseMetadata(bodyStream)
	if mErr != nil {
		return nil, mErr
	}

	switch metadata.cmdCode {
	case LoginCommandCode:
		return NewLoginCommand(*metadata, bodyStream, stream)
	case MessageCommandCode:
		return NewMessageCommand(*metadata, bodyStream)
	case CorrelationIDTestCommandCode:
		return NewCorrelationIDTestCommand(*metadata, bodyStream)
	default:
		return nil, ErrUnknownCommand
	}
}

func readFieldWithLength(stream io.Reader, fieldLen any) ([]byte, error) {
	var field []byte

	switch fieldLen.(type) {
	case uint16:

		tFieldLen := fieldLen.(uint16)
		err := binary.Read(stream, binary.BigEndian, &tFieldLen)
		if err != nil {
			return nil, err
		}

		field = make([]byte, tFieldLen)

	case uint32:

		tFieldLen := fieldLen.(uint32)
		err := binary.Read(stream, binary.BigEndian, &tFieldLen)
		if err != nil {
			return nil, err
		}

		field = make([]byte, tFieldLen)

	default:
		return nil, ErrUnsupportedLengthSize
	}

	_, err := io.ReadFull(stream, field)
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
