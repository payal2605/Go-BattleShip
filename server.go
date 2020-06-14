package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var playersOnline = make(map[string]*websocket.Conn)
var gamesLive = make(map[string]*Game)
var waiting []string
var gameId int = 1
var ongoing []int
var upgrader = websocket.Upgrader{}
var shipIndex []int = []int{5, 4, 3, 3, 2}

type Message struct {
	Task string `json:"task"`

	Message string `json:"message"`
	Turn    string `json:"turn"`
}

type Ship struct {
	start           string `json:"x"`
	size            int
	vertical        bool
	shipsCoordinate []string
	hit             int
}

type User struct {
	id                  string
	Ships               []Ship
	Coordinates         []string `json:"Coordinates"`
	AdjacentCoordinates []string
}

type Game struct {
	id      int
	players []User
}

func (g *Game) init(id int, player1id, player2id string) {
	g.id = id
	var p1, p2 User
	p1.init(player1id)
	p2.init(player2id)
	g.players = []User{p1, p2}
	gamesLive[player1id] = g
	gamesLive[player2id] = g
	coordinatesJson, err := json.Marshal(g.players[0].Coordinates)
	if err == nil {
		playersOnline[g.players[0].id].WriteJSON(Message{"sendcoordinates", string(coordinatesJson), "yourturn"})

	}
	coordinatesJson, err = json.Marshal(g.players[1].Coordinates)
	if err == nil {
		playersOnline[g.players[1].id].WriteJSON(Message{"sendcoordinates", string(coordinatesJson), "oppturn"})

	}
}

func (p *User) init(id string) {
	p.id = id

	p.CreateShips()
}
func (p *User) CreateShips() {

	for i := 0; i < len(shipIndex); i++ {
		var s Ship
		s.size = shipIndex[i]
		s.hit = 0
		fmt.Println(s.size)
		done := p.placeRandom(s)
		for !done {
			done = p.placeRandom(s)
		}
	}
}

func (p *User) placeRandom(s Ship) bool {
	tempCoordinates := make([]string, s.size)
	adjacentCoordinates := make([]string, 0)
	s.vertical = rand.Float64() < 0.5
	rand.Seed(time.Now().UnixNano())
	y := rand.Intn(10)
	x := rand.Intn(10)
	if s.vertical && x+s.size > 10 {
		return false
	} else if y+s.size > 10 {
		return false
	}
	if s.vertical {
		adjacentCoordinates = append(adjacentCoordinates, p.getAdjacent(x, y, s.size, "vertical")...)
	} else {
		adjacentCoordinates = append(adjacentCoordinates, p.getAdjacent(x, y, s.size, "horizontal")...)
	}
	ys := strconv.Itoa(y)
	xs := strconv.Itoa(x)
	s.start = xs + ys
	for i := 0; i < s.size; i++ {
		xs := strconv.Itoa(x)
		if s.vertical {
			tempCoordinates[i] = xs + ys
			if y > 0 {
				adjacentCoordinates = append(adjacentCoordinates, strconv.Itoa(x)+strconv.Itoa(y-1))
			}
			if y < 9 {
				adjacentCoordinates = append(adjacentCoordinates, strconv.Itoa(x)+strconv.Itoa(y+1))
			}
			x = x + 1
		} else {
			ys := strconv.Itoa(y)
			tempCoordinates[i] = xs + ys
			if x > 0 {
				adjacentCoordinates = append(adjacentCoordinates, strconv.Itoa(x-1)+strconv.Itoa(y))
			}
			if x < 9 {
				adjacentCoordinates = append(adjacentCoordinates, strconv.Itoa(x+1)+strconv.Itoa(y))
			}
			y = y + 1
		}

	}

	if !p.checkPresent(tempCoordinates, p.Coordinates) || !p.checkPresent(tempCoordinates, p.AdjacentCoordinates) {
		return false
	}

	p.Coordinates = append(p.Coordinates, tempCoordinates...)
	s.shipsCoordinate = append(s.shipsCoordinate, tempCoordinates...)
	p.AdjacentCoordinates = append(p.AdjacentCoordinates, adjacentCoordinates...)
	p.Ships = append(p.Ships, s)
	return true
}

func (p *User) getAdjacent(x, y, size int, direction string) []string {
	temp := make([]string, 0)

	if direction == "horizontal" {
		if y > 0 {
			back := strconv.Itoa(x) + strconv.Itoa(y-1)
			temp = append(temp, back)
			if x > 0 {
				backup := strconv.Itoa(x-1) + strconv.Itoa(y-1)
				temp = append(temp, backup)
			}
			if x < 9 {
				backdown := strconv.Itoa(x+1) + strconv.Itoa(y+1)
				temp = append(temp, backdown)
			}
		}
		if y+size-1 < 9 {
			back := strconv.Itoa(x) + strconv.Itoa(y+size)
			temp = append(temp, back)
			if x > 0 {
				backup := strconv.Itoa(x-1) + strconv.Itoa(y+size)
				temp = append(temp, backup)
			}
			if x < 9 {
				backdown := strconv.Itoa(x+1) + strconv.Itoa(y+size)
				temp = append(temp, backdown)
			}
		}
	} else {
		if x > 0 {
			back := strconv.Itoa(x-1) + strconv.Itoa(y)
			temp = append(temp, back)
			if y > 0 {
				backup := strconv.Itoa(x-1) + strconv.Itoa(y-1)
				temp = append(temp, backup)
			}
			if y < 9 {
				backdown := strconv.Itoa(x-1) + strconv.Itoa(y+1)
				temp = append(temp, backdown)
			}
		}
		if x+size-1 < 9 {
			back := strconv.Itoa(x+size) + strconv.Itoa(y)
			temp = append(temp, back)
			if y > 0 {
				backup := strconv.Itoa(x+size) + strconv.Itoa(y-1)
				temp = append(temp, backup)
			}
			if x < 9 {
				backdown := strconv.Itoa(x+size) + strconv.Itoa(y+1)
				temp = append(temp, backdown)
			}
		}

	}
	return temp
}

func (p *User) checkPresent(tempCoordinates []string, checkCoordinates []string) bool {

	for i := 0; i < len(tempCoordinates); i++ {
		for j := 0; j < len(checkCoordinates); j++ {
			if tempCoordinates[i] == checkCoordinates[j] {
				return false
			}
		}

	}

	return true
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err, "err")
			return
		}

		go HandleClient(conn)
	})

	fmt.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func HandleClient(c *websocket.Conn) {
	uuidv4, _ := uuid.NewV4()
	id := uuidv4.String()
	playersOnline[id] = c
	fmt.Printf("%s ID %s connected.\n", time.Now().String(), id)
	waiting = append(waiting, id)
	checkPlayers()

	for {
		var msg Message

		err := c.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err, "err")
			handleMessage(id, Message{"disconnect", "", ""})
			return
		}
		if msg.Task != "" {
			handleMessage(id, msg)
		}
	}
}

func handleMessage(id string, msg Message) {
	var me, opp User

	if gamesLive[id].players[0].id == id {
		me = gamesLive[id].players[0]
		opp = gamesLive[id].players[1]

	} else {
		me = gamesLive[id].players[1]
		opp = gamesLive[id].players[0]
	}
	if msg.Task == "disconnect" {
		playersOnline[opp.id].WriteJSON(Message{"disconnect", "", ""})
	}

	if msg.Task == "win" {
		playersOnline[me.id].WriteJSON(Message{"win", "", ""})
		playersOnline[opp.id].WriteJSON(Message{"lose", "", ""})
	}
	if msg.Task == "shot" {
		me.checkShot(&opp, msg.Message)
	}

}
func (me *User) checkShot(opp *User, message string) {
	var flag = 0
	for i := 0; i < len(opp.Ships); i++ {
		for j := 0; j < len(opp.Ships[i].shipsCoordinate); j++ {
			if message == opp.Ships[i].shipsCoordinate[j] {
				opp.Ships[i].hit++
				playersOnline[me.id].WriteJSON(Message{"hit", message, "me"})
				playersOnline[opp.id].WriteJSON(Message{"hit", message, "you"})
				if opp.Ships[i].hit == len(opp.Ships[i].shipsCoordinate) {
					coordinatesJson, err := json.Marshal(opp.Ships[i].shipsCoordinate)
					if err == nil {
						playersOnline[me.id].WriteJSON(Message{"hitShip", string(coordinatesJson), ""})
					}
				}
				flag = 1
				break
			}
		}
	}

	if flag == 0 {
		playersOnline[me.id].WriteJSON(Message{"miss", message, "me"})
		playersOnline[opp.id].WriteJSON(Message{"miss", message, "you"})
	}
}

func checkPlayers() {
	var player1, player2 string
	if len(waiting) >= 2 {
		player1 = waiting[0]
		player2 = waiting[1]

		waiting = waiting[2:]
		var game Game
		fmt.Println(gameId)
		game.init(gameId, player1, player2)
		gameId++
	}
}
