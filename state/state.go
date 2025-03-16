package state

import "errors"

type State struct {
	loggedUsers map[string]bool
	messages    map[string](chan string)
}

var (
	ErrUserAlreadyOnline = errors.New("user already online")
)

func NewState() *State {
	return &State{
		loggedUsers: map[string]bool{},
		messages:    map[string](chan string){},
	}
}

func (s *State) UserExists(username string) bool {
	_, ok := s.loggedUsers[username]
	return ok
}

func (s *State) UserIsOnline(username string) bool {
	return s.loggedUsers[username]
}

func (s *State) Login(username string) error {

	if s.UserIsOnline(username) {
		return ErrUserAlreadyOnline
	}

	s.loggedUsers[username] = true
	return nil
}

func (s *State) Logout(username string) {
	s.loggedUsers[username] = false
}
