package commands

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

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
			body: "\x00\x00\x00\x11\x01\x00\x01\x00\x00\x00\x01\x00\x08TestUser",
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
			name: "happy path: correct correlationIDTest packet gets parsed",
			body: "\x00\x00\x00\x07\x01\x00\x09\x00\x00\x00\x0A",
			wantRes: &CorrelationIDTestCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       9,
					correlationId: 10,
				},
			},
			wantErr: nil,
		},

		{
			name: "happy path: correct message packet gets parsed",
			body: "\x00\x00\x00\x1E\x01\x00\x02\x00\x00\x00\x01\x00\x03msg\x00\x03usr\x00\x03rec\x18\x16\x68\x7E\xC0\x57\x00\x00",
			wantRes: &MessageCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       2,
					correlationId: 1,
				},
				message:   "msg",
				from:      "usr",
				to:        "rec",
				timestamp: time.Unix(1735689600, 0),
			},
			wantErr: nil,
		},
		{
			name:    "error: malformed command, length field too short",
			body:    "\x01\x00",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, incorrect length",
			body:    "\x00\x00\x00\x11\x01",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, missing metadata",
			body:    "\x00\x00\x00\x00",
			wantRes: nil,
			wantErr: errors.New("EOF"),
		},
		{
			name:    "error: malformed metadata, command code too short",
			body:    "\x00\x00\x00\x02\x01\x00",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, correlationId too short",
			body:    "\x00\x00\x00\x06\x01\x00\x01\x00\x00\x00",
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: unknown command",
			body:    "\x00\x00\x00\x07\x01\x00\x99\x00\x00\x00\x01",
			wantRes: nil,
			wantErr: ErrUnknownCommand,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			buf := bufio.NewReader(bytes.NewBuffer([]byte(tt.body)))

			res, err := ParseCommand(buf)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
