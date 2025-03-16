package commands

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ParseCommand(t *testing.T) {
	mockLoginStream := generateStream("\x00\x00\x00\x11\x01\x00\x01\x00\x00\x00\x01\x00\x08TestUser")
	tests := []struct {
		name    string
		stream  net.Conn
		wantRes Command
		wantErr error
	}{
		{
			name:   "happy path: correct login packet gets parsed",
			stream: mockLoginStream,
			wantRes: &LoginCommand{
				metadata: Metadata{
					version:       1,
					cmdCode:       1,
					correlationId: 1,
				},
				username: "TestUser",
				conn:     mockLoginStream,
			},
			wantErr: nil,
		},
		{
			name:   "happy path: correct correlationIDTest packet gets parsed",
			stream: generateStream("\x00\x00\x00\x07\x01\x00\x09\x00\x00\x00\x0A"),
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
			name:   "happy path: correct message packet gets parsed",
			stream: generateStream("\x00\x00\x00\x1E\x01\x00\x02\x00\x00\x00\x01\x00\x03msg\x00\x03usr\x00\x03rec\x18\x16\x68\x7E\xC0\x57\x00\x00"),
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
			stream:  generateStream("\x01\x00"),
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, incorrect length",
			stream:  generateStream("\x00\x00\x00\x11\x01"),
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, missing metadata",
			stream:  generateStream("\x00\x00\x00\x00"),
			wantRes: nil,
			wantErr: errors.New("EOF"),
		},
		{
			name:    "error: malformed metadata, command code too short",
			stream:  generateStream("\x00\x00\x00\x02\x01\x00"),
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: malformed command, correlationId too short",
			stream:  generateStream("\x00\x00\x00\x06\x01\x00\x01\x00\x00\x00"),
			wantRes: nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "error: unknown command",
			stream:  generateStream("\x00\x00\x00\x07\x01\x00\x99\x00\x00\x00\x01"),
			wantRes: nil,
			wantErr: ErrUnknownCommand,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := ParseCommand(tt.stream)

			assert.Equal(t, tt.wantRes, res)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func generateStream(body string) net.Conn {
	server, client := net.Pipe()
	go func() {
		_, _ = client.Write([]byte(body))
		client.Close()
	}()
	return server
}
