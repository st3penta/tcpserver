package commands

import (
	"bufio"
	"bytes"
	"tcpserver/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCorrelationIDTestCommand(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantRes *CorrelationIDTestCommand
		wantErr error
	}{
		{
			name: "happy path: correct correlationIDTest packet gets parsed",
			body: "",
			wantRes: &CorrelationIDTestCommand{
				metadata: Metadata{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := bufio.NewReader(bytes.NewBuffer([]byte(tt.body)))

			res, err := NewCorrelationIDTestCommand(Metadata{}, buf)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_CorrelationIDTestCommand_Process(t *testing.T) {
	tests := []struct {
		name    string
		cc      *CorrelationIDTestCommand
		wantRes *Response
		wantErr error
	}{
		{
			name: "happy path: CorrelationID Test command gets processed",
			cc: &CorrelationIDTestCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       9,
					correlationId: 10,
				},
			},
			wantRes: &Response{
				version:       1,
				correlationID: 10,
				statusCode:    1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := tt.cc.Process(state.NewState())

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
