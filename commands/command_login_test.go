package commands

import (
	"bufio"
	"bytes"
	"io"
	"tcpserver/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLoginCommand(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantRes *LoginCommand
		wantErr error
	}{
		{
			name: "happy path: correct login packet gets parsed",
			body: "\x00\x08TestUser",
			wantRes: &LoginCommand{
				metadata: Metadata{},
				username: "TestUser",
			},
			wantErr: nil,
		},
		{
			name:    "error: malformed command, length field too short",
			body:    "\x01",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, username length incorrect",
			body:    "\x00\x08short",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := bufio.NewReader(bytes.NewBuffer([]byte(tt.body)))

			res, err := NewLoginCommand(Metadata{}, buf)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_LoginCommand_Process(t *testing.T) {
	tests := []struct {
		name    string
		lc      *LoginCommand
		wantRes *Response
		wantErr error
	}{
		{
			name: "happy path: login command gets processed",
			lc: &LoginCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       LoginCommandCode,
					correlationId: 1,
				},
				username: "TestUser",
			},
			wantRes: &Response{
				version:       1,
				correlationID: 1,
				statusCode:    1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := tt.lc.Process(state.NewState())

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
