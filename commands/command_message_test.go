package commands

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"tcpserver/state"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewMessageCommand(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantRes *MessageCommand
		wantErr error
	}{
		{
			name: "happy path: correct message packet gets parsed",
			body: "\x00\x03msg\x00\x03usr\x00\x03rec\x18\x16\x68\x7E\xC0\x57\x00\x00",
			wantRes: &MessageCommand{
				metadata:  Metadata{},
				message:   "msg",
				from:      "usr",
				to:        "rec",
				timestamp: time.Unix(1735689600, 0),
			},
			wantErr: nil,
		},
		{
			name:    "error: malformed command, message length field too short",
			body:    "\x01",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, message length incorrect",
			body:    "\x00\x08short",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, from length field too short",
			body:    "\x00\x03msg\x00",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, from length incorrect",
			body:    "\x00\x03msg\x00\x08short",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, to length field too short",
			body:    "\x00\x03msg\x00\x03usr\x00",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, to length incorrect",
			body:    "\x00\x03msg\x00\x03usr\x00\x08short",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, timestamp field too short",
			body:    "\x00\x03msg\x00\x03usr\x00\x03rec\x18\x16\x68\x7E\xC0\x57",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := bufio.NewReader(bytes.NewBuffer([]byte(tt.body)))

			res, err := NewMessageCommand(Metadata{}, buf)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_MessageCommand_Process(t *testing.T) {
	mockConn := net.TCPConn{}
	tests := []struct {
		name    string
		lc      *MessageCommand
		state   State
		wantRes *Response
		wantErr error
	}{
		{
			name: "happy path: message command gets processed",
			lc: &MessageCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       MessageCommandCode,
					correlationId: 1,
				},
				from:      "sender",
				to:        "recipient",
				timestamp: time.Time{},
				message:   "message",
			},
			state: func() *state.State {
				s := state.NewState()

				_ = s.Login(&mockConn, "recipient")
				return s
			}(),
			wantRes: &Response{
				version:       1,
				correlationID: 1,
				statusCode:    1,
			},
			wantErr: nil,
		},
		{
			name: "error: recipient doesn't exist",
			lc: &MessageCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       MessageCommandCode,
					correlationId: 1,
				},
				from:      "sender",
				to:        "recipient",
				timestamp: time.Time{},
				message:   "message",
			},
			state:   state.NewState(),
			wantRes: nil,
			wantErr: state.ErrRecipientNotExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := tt.lc.Process(tt.state)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
