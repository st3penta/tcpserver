package commands

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"tcpserver/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLoginCommand(t *testing.T) {
	mockConn := net.TCPConn{}

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
				conn:     &mockConn,
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

			res, err := NewLoginCommand(Metadata{}, buf, &mockConn)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_LoginCommand_Process(t *testing.T) {
	mockConn1 := net.TCPConn{}
	mockConn2 := net.TCPConn{}
	tests := []struct {
		name      string
		lc        *LoginCommand
		state     state.State
		wantRes   *Response
		wantState state.State
		wantErr   error
	}{
		{
			name: "happy path: login command gets processed",
			lc: &LoginCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       LoginCommandCode,
					correlationId: 1,
				},
				username: "user2",
				conn:     &mockConn2,
			},
			state: state.State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
				Connections: map[net.Conn]string{
					&mockConn1: "user1",
				},
			},
			wantRes: &Response{
				version:       1,
				correlationID: 1,
				statusCode:    1,
			},
			wantState: state.State{
				LoggedUsers: map[string]bool{
					"user1": true,
					"user2": true,
				},
				Connections: map[net.Conn]string{
					&mockConn1: "user1",
					&mockConn2: "user2",
				},
			},
			wantErr: nil,
		},
		{
			name: "error: user already online",
			lc: &LoginCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       LoginCommandCode,
					correlationId: 1,
				},
				username: "user1",
			},
			state: state.State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
			},
			wantRes: nil,
			wantState: state.State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
			},
			wantErr: state.ErrUserAlreadyOnline,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := tt.lc.Process(&tt.state)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantState, tt.state)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
