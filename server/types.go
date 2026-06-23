package server

// Player is one joined player from an ASE full query response.
type Player struct {
	Name  string
	Score int
	Ping  int
}
