package state

import (
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	MessageQueueMaxSize = 100
)

var (
	ErrUserAlreadyOnline  = errors.New("user already online")
	ErrRecipientNotExists = errors.New("recipent doesn't exist")
)

type State struct {
	Connections map[io.Reader]string
	LoggedUsers map[string]bool
	Messages    map[string](chan Message)
}

func NewState() *State {
	return &State{
		LoggedUsers: map[string]bool{},
		Messages:    map[string](chan Message){},
		Connections: map[io.Reader]string{},
	}
}

func (s *State) Login(conn io.Reader, username string) error {

	if s.userIsOnline(username) {
		return ErrUserAlreadyOnline
	}

	s.LoggedUsers[username] = true
	s.Connections[conn] = username

	return nil
}

func (s *State) Logout(conn io.Reader) {
	s.LoggedUsers[s.Connections[conn]] = false
	delete(s.Connections, conn)
}

func (s *State) EnqueueMessage(from string, to string, timestamp time.Time, message string) error {

	if !s.userExists(to) {
		return ErrRecipientNotExists
	}

	_, ok := s.Messages[to]
	if !ok {
		s.Messages[to] = make(chan Message, MessageQueueMaxSize)
	}

	s.Messages[to] <- Message{
		From:      from,
		Timestamp: timestamp,
		Payload:   message,
	}

	return nil
}

func (s *State) userExists(username string) bool {
	_, ok := s.LoggedUsers[username]
	return ok
}

func (s *State) userIsOnline(username string) bool {
	return s.LoggedUsers[username]
}

type Message struct {
	From      string
	Timestamp time.Time
	Payload   string
}

func (m *Message) Print() {
	fmt.Println("-----")
	fmt.Println("Message")
	fmt.Printf("\tfrom: %s\n", m.From)
	fmt.Printf("\ttime: %s\n", m.Timestamp.String())
	fmt.Printf("\tpayload: %s\n", m.Payload)
	fmt.Println("-----")
}
