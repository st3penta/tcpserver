package commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	MessageCommandCode uint16 = 0x02
)

type MessageCommand struct {
	metadata  Metadata
	message   string
	from      string
	to        string
	timestamp time.Time
}

func NewMessageCommand(
	metadata Metadata,
	stream io.Reader,
) (*MessageCommand, error) {

	var mLen uint16
	message, mErr := readFieldWithLength(stream, mLen)
	if mErr != nil {
		return nil, mErr
	}

	var fLen uint16
	from, fErr := readFieldWithLength(stream, fLen)
	if fErr != nil {
		return nil, fErr
	}

	var tLen uint16
	to, tErr := readFieldWithLength(stream, tLen)
	if tErr != nil {
		return nil, tErr
	}

	var timestamp int64
	err := binary.Read(stream, binary.BigEndian, &timestamp)
	if err != nil {
		return nil, err
	}

	mc := &MessageCommand{
		metadata:  metadata,
		message:   string(message),
		from:      string(from),
		to:        string(to),
		timestamp: time.Unix(0, timestamp),
	}

	mc.print()

	return mc, nil
}

func (mc *MessageCommand) Process(state State) (*Response, error) {

	// TODO actual message processing

	return &Response{
		version:       mc.metadata.version,
		correlationID: mc.metadata.correlationId,
		statusCode:    ResponseStatusCodeOK,
	}, nil
}

func (mc *MessageCommand) print() {
	fmt.Println("-----")
	fmt.Println("Message")
	fmt.Printf("\tversion: %d\n", mc.metadata.version)
	fmt.Printf("\tcorrelationId: %d\n", mc.metadata.correlationId)
	fmt.Printf("\tmessage: %s\n", mc.message)
	fmt.Printf("\tfrom: %s\n", mc.from)
	fmt.Printf("\tto: %s\n", mc.to)
	fmt.Printf("\ttime: %s\n", mc.timestamp.String())
	fmt.Println("-----")
}
