package state

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_State_Login(t *testing.T) {
	mockConn := net.TCPConn{}

	tests := []struct {
		name      string
		state     *State
		conn      net.Conn
		username  string
		wantState *State
		wantErr   error
	}{
		{
			name:     "happy path",
			state:    NewState(),
			conn:     &mockConn,
			username: "user1",
			wantState: &State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn: "user1",
				},
				Interrupts: map[string]chan bool{},
			},
			wantErr: nil,
		},
		{
			name: "error: user already online",
			state: &State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn: "user1",
				},
				Interrupts: map[string]chan bool{},
			},
			conn:     &mockConn,
			username: "user1",
			wantState: &State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn: "user1",
				},
				Interrupts: map[string]chan bool{},
			},
			wantErr: ErrUserAlreadyOnline,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.state.Login(tt.conn, tt.username)

			assert.Equal(t, *tt.wantState, *tt.state)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_State_Logout(t *testing.T) {
	mockConn := net.TCPConn{}

	tests := []struct {
		name      string
		state     *State
		conn      net.Conn
		wantState *State
	}{
		{
			name: "happy path",
			state: &State{
				LoggedUsers: map[string]bool{
					"user1": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn: "user1",
				},
			},
			conn: &mockConn,
			wantState: &State{
				LoggedUsers: map[string]bool{
					"user1": false,
				},
				Messages:    map[string]chan Message{},
				Connections: map[net.Conn]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.state.Logout(tt.conn)

			assert.Equal(t, *tt.wantState, *tt.state)
		})
	}
}

func Test_State_EnqueueMessage(t *testing.T) {
	mockConn1 := net.TCPConn{}
	mockConn2 := net.TCPConn{}

	tests := []struct {
		name         string
		state        *State
		from         string
		to           string
		timestamp    time.Time
		message      string
		wantMessages []Message
		wantErr      error
	}{
		{
			name: "happy path",
			state: &State{
				LoggedUsers: map[string]bool{
					"sender":    true,
					"recipient": true,
				},
				Messages: map[string]chan Message{
					"recipient": make(chan Message, MessageQueueMaxSize),
				},
				Connections: map[net.Conn]string{
					&mockConn1: "sender",
					&mockConn2: "recipient",
				},
				Interrupts: map[string]chan bool{
					"sender":    make(chan bool, 1),
					"recipient": make(chan bool, 1),
				},
			},
			from:      "sender",
			to:        "recipient",
			timestamp: time.Time{},
			message:   "message",
			wantMessages: []Message{
				{
					From:      "sender",
					Timestamp: time.Time{},
					Payload:   "message",
				},
			},
			wantErr: nil,
		},
		{
			name: "happy path, channel must be created",
			state: &State{
				LoggedUsers: map[string]bool{
					"sender":    true,
					"recipient": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn1: "sender",
					&mockConn2: "recipient",
				},
				Interrupts: map[string]chan bool{
					"sender":    make(chan bool, 1),
					"recipient": make(chan bool, 1),
				},
			},
			from:      "sender",
			to:        "recipient",
			timestamp: time.Time{},
			message:   "message",
			wantMessages: []Message{
				{
					From:      "sender",
					Timestamp: time.Time{},
					Payload:   "message",
				},
			},
			wantErr: nil,
		},
		{
			name: "happy path, recipient is offline and msg gets enqueued",
			state: &State{
				LoggedUsers: map[string]bool{
					"sender":    true,
					"recipient": false,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn1: "sender",
					&mockConn2: "recipient",
				},
				Interrupts: map[string]chan bool{
					"sender": make(chan bool, 1),
				},
			},
			from:      "sender",
			to:        "recipient",
			timestamp: time.Time{},
			message:   "message",
			wantMessages: []Message{
				{
					From:      "sender",
					Timestamp: time.Time{},
					Payload:   "message",
				},
			},
			wantErr: nil,
		},
		{
			name: "error: recipient doesn't exist",
			state: &State{
				LoggedUsers: map[string]bool{
					"sender": true,
				},
				Messages: map[string]chan Message{},
				Connections: map[net.Conn]string{
					&mockConn1: "sender",
					&mockConn2: "recipient",
				},
				Interrupts: map[string]chan bool{
					"sender": make(chan bool, 1),
				},
			},
			from:         "sender",
			to:           "recipient",
			timestamp:    time.Time{},
			message:      "message",
			wantMessages: []Message{},
			wantErr:      ErrRecipientNotExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			received := []Message{}
			go collectReceived(tt.state.Messages, &received)

			err := tt.state.EnqueueMessage(tt.from, tt.to, tt.timestamp, tt.message)

			time.Sleep(20 * time.Millisecond)
			assert.Equal(t, tt.wantMessages, received)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func collectReceived(msgsMap map[string]chan Message, received *[]Message) {

	for {
		msgs := msgsMap["recipient"]
		if msgs == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		msg := <-msgs
		*received = append(*received, msg)
		msg.Print()
	}
}
