package commands

import (
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
			body: "\x01\x00\x01\x00\x00\x00\x01\x00\x08TestUser",
			wantRes: &LoginCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       1,
					correlationId: 1,
				},
				username: "TestUser",
			},
			wantErr: nil,
		},
		{
			name:    "malformed packet: incomplete metadata",
			body:    "\x01\x00\x01\x00",
			wantRes: nil,
			wantErr: ErrMalformedMetadata,
		},
		{
			name:    "malformed packet: command is too short",
			body:    "\x01\x00\x01\x00\x00\x00\x01\x00",
			wantRes: nil,
			wantErr: ErrLoginCommandTooShort,
		},
		{
			name:    "happy path: invalid username length",
			body:    "\x01\x00\x01\x00\x00\x00\x01\x00\x08short",
			wantRes: nil,
			wantErr: ErrInvalidUsernameLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := NewLoginCommand([]byte(tt.body))

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
					cmdCode:       1,
					correlationId: 1,
				},
				username: "TestUser",
			},
			wantRes: &Response{
				responseLength: 9,
				Metadata: Metadata{
					version:       1,
					cmdCode:       3,
					correlationId: 1,
				},
				statusCode: 1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := tt.lc.Process()

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
