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
	fmt.Println(fmt.Sprintf("   version: %d", cc.metadata.version))
	fmt.Println(fmt.Sprintf("   correlationId: %d", cc.metadata.correlationId))
	fmt.Println("-----")
}
