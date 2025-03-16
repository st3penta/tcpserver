package commands

import (
	"fmt"
	"io"
)

var (
	CorrelationIDTestResponseLength uint32 = 0x0009
	CorrelationIDTestCommandCode    uint16 = 0x09
)

type CorrelationIDTestCommand struct {
	metadata Metadata
}

func NewCorrelationIDTestCommand(
	metadata Metadata,
	stream io.Reader,
) (*CorrelationIDTestCommand, error) {

	cc := &CorrelationIDTestCommand{
		metadata: metadata,
	}

	cc.print()

	return cc, nil
}

func (cc *CorrelationIDTestCommand) Process() (*Response, error) {
	return &Response{
		version:       cc.metadata.version,
		correlationID: cc.metadata.correlationId,
		statusCode:    ResponseStatusCodeOK,
	}, nil
}

func (cc *CorrelationIDTestCommand) print() {
	fmt.Println("-----")
	fmt.Println("CorrelationID Test")
	fmt.Printf("   version: %d\n", cc.metadata.version)
	fmt.Printf("   correlationId: %d\n", cc.metadata.correlationId)
	fmt.Println("-----")
}
