package state

type State struct {
	loggedUsers map[string]bool
	messages    map[string](chan string)
}

func NewState() *State {
	return &State{
		loggedUsers: map[string]bool{},
		messages:    map[string](chan string){},
	}
}
