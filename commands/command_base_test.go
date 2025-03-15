package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantRes Command
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
			name:    "error: malformed command",
			body:    "\x01\x00",
			wantRes: nil,
			wantErr: ErrMalformedCommand,
		},
		{
			name:    "error: unknown command",
			body:    "\x01\x00\x99\x00",
			wantRes: nil,
			wantErr: ErrUnknownCommand,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := ParseCommand([]byte(tt.body))

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
