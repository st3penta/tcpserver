package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Response_Write(t *testing.T) {
	tests := []struct {
		name       string
		response   *Response
		wantOutput string
		wantErr    error
	}{
		{
			name: "happy path: login command gets processed",
			response: &Response{
				version:       1,
				correlationID: 1,
				statusCode:    1,
			},
			wantOutput: "\x00\x00\x00\x09\x01\x00\x03\x00\x00\x00\x01\x00\x01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer

			err := tt.response.Write(&buf)

			assert.Equal(t, []byte(tt.wantOutput), buf.Bytes())
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
