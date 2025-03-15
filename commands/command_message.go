package commands

import (
	"fmt"
	"io"
)

type MessageCommand struct {
	Metadata
	CorrelationId uint32
	Message       string
	From          string
	To            string
	Time          string
}

type MessageResponse = BaseResponse

func NewMessageCommand(
	version byte,
	cmdCode uint16,
	correlationId uint32,
	message string,
	from string,
	to string,
	time string,
) *MessageCommand {
	return &MessageCommand{
		Metadata: Metadata{
			version: version,
			cmdCode: cmdCode,
		},
		CorrelationId: correlationId,
		Message:       message,
		From:          from,
		To:            to,
		Time:          time,
	}
}

func (lc *MessageCommand) Process(out io.Writer) {
	resp := &MessageResponse{
		Metadata: Metadata{
			version: lc.version,
			cmdCode: 3,
		},
		CorrelationId: lc.CorrelationId,
		StatusCode:    1,
	}

	resp.Write(out)
}

func (mc *MessageCommand) Print() {
	fmt.Println("-----")
	fmt.Println("Message")
	fmt.Println(fmt.Sprintf("version: %d", mc.version))
	fmt.Println(fmt.Sprintf("correlationId: %d", mc.CorrelationId))
	fmt.Println(fmt.Sprintf("message: %s", mc.Message))
	fmt.Println(fmt.Sprintf("from: %s", mc.From))
	fmt.Println(fmt.Sprintf("to: %s", mc.To))
	fmt.Println(fmt.Sprintf("time: %s", mc.Time))
	fmt.Println("-----")
}
