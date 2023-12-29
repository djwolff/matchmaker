package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/djwolff/matchmaker/models/forms"
	"github.com/djwolff/matchmaker/models/mm"
	"github.com/djwolff/matchmaker/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type MatchmakingServer struct {
	players     map[string]*mm.Player
	matchQueue  map[string][]*mm.Player // matchQueue containts buckets of gameModes.
	matchBuffer []*mm.Match
	mutex       sync.Mutex // prevent race conditions
}

func NewMatchmakingServer() *MatchmakingServer {
	s := &MatchmakingServer{
		players:     make(map[string]*mm.Player),
		matchQueue:  make(map[string][]*mm.Player),
		matchBuffer: make([]*mm.Match, 0),
	}
	// heap.Init(&s.matchQueue[0])
	return s
}

// Only focus on league of legends for now
func (s *MatchmakingServer) MatchMake(c *gin.Context, gormDB *gorm.DB, videogame string) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, "User is not Logged in")
		return
	}

	userStruct, ok := user.(*token.JWTUser)
	if !ok {
		c.JSON(http.StatusUnauthorized, "Failed to cast user struct")
		return
	}
	// curUser, err := json.Marshal(&user)
	form := new(forms.LeaguePlayerForm)
	if err := c.Bind(form); err != nil {
		c.JSON(http.StatusBadRequest, "Form is not for League of Legends")
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"Validation errors": validationErrors})
		return
	}

	// Put player in match making pool and find matcher (Nakama and https://github.com/didopimentel/matchmaker inspired)
	// Will run match making service on local -> cloud
	// will require open socket connection with each user
	// Make a request ticket, put ticket into pool of available players.
	// Pool by Region
	// https://github.com/heroiclabs/nakama/blob/master/server/matchmaker.go
	s.addPlayer(c.Writer, c.Request, userStruct, form)
}

func (s *MatchmakingServer) addPlayer(w http.ResponseWriter, r *http.Request, user *token.JWTUser, form *forms.LeaguePlayerForm) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player := &mm.Player{
		ID:           user.ID, // ID will be user ID, enforce user queue once
		GameMode:     form.GameMode,
		SeekingRole:  form.SeekingRole,
		OfferingRole: form.OfferingRole,
		Conn:         conn,
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.players[player.ID] = player

	// Add player to game mode's PQ
	s.matchQueue[form.GameMode] = append(s.matchQueue[form.GameMode], player)
}

func (s *MatchmakingServer) ContinuousMatchmaking() {
	go func() {
		for {
			time.Sleep(5 * time.Second) // Adjust the sleep duration as needed
			s.tryMatch()
		}
	}()
}
func (s *MatchmakingServer) tryMatch() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// find matches within matchQueue of players
	// for each bucket, go routine on finding matches within the bucket
	for gameMode, players := range s.matchQueue {
		// TODO: go routine on this search
		fmt.Println("finding matches for gamemode: ", gameMode)

		sort.Slice(players, func(i, j int) bool {
			return players[i].OfferingRole < players[j].OfferingRole
		})

		// Find the list of maximum exclusive matches and unmatched players
		matches, unmatchedPlayers := maxExclusiveMatches(players)
		s.matchBuffer = append(s.matchBuffer, matches...)

		for _, match := range matches {
			// Notify players about the match using websockets
			match.Player1.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You are matched with %s", match.Player2.ID)))
			match.Player2.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You are matched with %s", match.Player1.ID)))

			fmt.Printf("Matched players: %s and %s\n", match.Player1.ID, match.Player2.ID)

			// Close the websockets after the match
			match.Player1.Conn.Close()
			match.Player2.Conn.Close()

			// Link each player the other player's discord?
			// discord.com/users/{user_id}
			// https://discord.com/developers/docs/resources/user#get-user

		}
		players = unmatchedPlayers
	}
}

func maxExclusiveMatches(players []*mm.Player) (matches []*mm.Match, unmatchedPlayers []*mm.Player) {
	// Create a bipartite graph representation
	graph := make(map[string][]string)
	for _, player := range players {
		graph[player.ID] = append(graph[player.ID], player.SeekingRole)
	}

	// Perform Hopcroft-Karp algorithm
	matching := make(map[string]string)
	for _, player := range players {
		if matching[player.ID] == "" {
			visited := make(map[string]bool)
			augmentPath(graph, player.ID, visited, matching)
		}
	}

	// Build the list of matches and unmatched players
	for _, player := range players {
		if matching[player.ID] != "" {
			matches = append(matches, &mm.Match{Player1: player, Player2: &mm.Player{ID: matching[player.ID]}})
			matching[player.ID] = "" // Mark the player as matched
		} else {
			unmatchedPlayers = append(unmatchedPlayers, player)
		}
	}

	return matches, unmatchedPlayers
}

func augmentPath(graph map[string][]string, playerID string, visited map[string]bool, matching map[string]string) bool {
	for _, role := range graph[playerID] {
		if !visited[role] {
			visited[role] = true
			if matching[role] == "" || augmentPath(graph, matching[role], visited, matching) {
				matching[role] = playerID
				return true
			}
		}
	}
	return false
}
