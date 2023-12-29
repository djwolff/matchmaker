package mm

import "github.com/gorilla/websocket"

type Player struct {
	ID string
	// Skill        int
	GameMode     string
	SeekingRole  string
	OfferingRole string
	Conn         *websocket.Conn // each player requires a conn to connect to pool
	Index        int
}

// This match is only for one on one
type Match struct {
	MatchID  uint
	GameMode *string
	Player1  *Player
	Player2  *Player
}
