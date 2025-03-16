package state

import (
	"errors"
	"io"
)

type State struct {
	Connections map[io.Reader]string
	LoggedUsers map[string]bool
	Messages    map[string](chan string)
}

var (
	ErrUserAlreadyOnline = errors.New("user already online")
)

func NewState() *State {
	return &State{
		LoggedUsers: map[string]bool{},
		Messages:    map[string](chan string){},
		Connections: map[io.Reader]string{},
	}
}

func (s *State) UserExists(username string) bool {
	_, ok := s.LoggedUsers[username]
	return ok
}

func (s *State) UserIsOnline(username string) bool {
	return s.LoggedUsers[username]
}

func (s *State) Login(conn io.Reader, username string) error {

	if s.UserIsOnline(username) {
		return ErrUserAlreadyOnline
	}

	s.LoggedUsers[username] = true
	s.Connections[conn] = username

	return nil
}

func (s *State) Logout(username string) {
	s.LoggedUsers[username] = false
}
