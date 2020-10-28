package broker

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/knodesec/flexitty"
)

var Sessions []*Session

type Session struct {
	UUID uuid.UUID
	WS   []*websocket.Conn
	TTY  *flexitty.TTY
}

func (s *Session) AddWS(c *websocket.Conn) {
	s.WS = append(s.WS, c)

	// Start a goroutine for handling input from the termjs window in the browser
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ERROR: Couldnt read websocket ", err)
				break
			}
			s.TTY.Write(message)
		}
	}()

	log.Printf("DEBUG: AddWS - Added websocket to session %s\n", s.UUID.String())
	log.Printf("DEBUG: AddWS - Session has %d sockets\n", len(s.WS))
}

func (s Session) Broadcast(msgtype int, data []byte) error {
	log.Printf("DEBUG: s.Broadcast - NUmber of sessions: %d\n", len(Sessions))
	log.Printf("DEBUG: s.Broadcast - Starting..\n")
	log.Printf("DEBUG: s.Broadcast - Session has %d websockets\n", len(s.WS))
	for i, wsock := range s.WS {
		log.Printf("DEBUG: s.Broadcast - [%d] - %s\n", i, wsock.UnderlyingConn().RemoteAddr().String())
		wsock.WriteMessage(msgtype, data)
	}
	log.Printf("DEBUG: s.Broadcast - broadcast\n")
	return nil
}

func (s *Session) StartTTYReader() {
	log.Printf("DEBUG: s.TTYReader - Starting goroutine\n")
	go func() {
		for {
			log.Printf("DEBUG:TTYReader - Top of the for loop\n")
			var data []byte
			data, err := s.TTY.Read()
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("DEBUG: s.TTYReader - read some data\n")
			s.Broadcast(websocket.TextMessage, data)
		}
	}()
}
func NewSession() string {
	sesh := Session{
		UUID: uuid.New(),
	}

	newTTY, err := flexitty.New("bash", []string{})
	if err != nil {
		panic(err)
	}

	sesh.TTY = newTTY
	sesh.StartTTYReader()

	Sessions = append(Sessions, &sesh)

	return sesh.UUID.String()
}

func SessionExists(uuidString string) bool {
	for _, s := range Sessions {
		if s.UUID.String() == uuidString {
			return true
		}
	}
	return false
}

func AddWS(uuid string, c *websocket.Conn) error {

	for _, s := range Sessions {
		if s.UUID.String() == uuid {
			s.AddWS(c)
			return nil
		}
	}
	return fmt.Errorf("adding websocket to broker failed, uuid %s did not exist", uuid)
}

