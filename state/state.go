package state

import (
	"errors"
	"fmt"
	"net"
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
	Connections map[net.Conn]string
	LoggedUsers map[string]bool
	Messages    map[string]chan Message
	Interrupts  map[string]chan bool
}

func NewState() *State {
	return &State{
		Connections: map[net.Conn]string{},
		LoggedUsers: map[string]bool{},
		Messages:    map[string](chan Message){},
		Interrupts:  map[string]chan bool{},
	}
}

func (s *State) Login(conn net.Conn, username string) error {

	if s.userIsOnline(username) {
		return ErrUserAlreadyOnline
	}

	s.LoggedUsers[username] = true
	s.Connections[conn] = username

	return nil
}

func (s *State) Logout(conn net.Conn) {
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

	_, ok = s.Interrupts[to]
	if !ok {
		s.Interrupts[to] = make(chan bool, 1)
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
