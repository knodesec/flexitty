package broker

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/knodesec/flexitty"
)

var Manager = ManagerStruct{}

type ManagerStruct struct {
	Sessions []Session
}

type Session struct {
	UUID uuid.UUID
	WS   *websocket.Conn
	TTY  *flexitty.TTY
}

func (s *Session) AddWS(c *websocket.Conn) {
	s.WS = c
	go func() {
		for {
			_, message, err := s.WS.ReadMessage()
			if err != nil {
				log.Println("WS read err: ", err)
				break
			}
			//log.Printf("WS recv: %s\n", message)
			s.TTY.InputChan <- message
		}
	}()

	go func() {
		for data := range s.TTY.OutputChan {
			//log.Printf("Data from the TTY: %v\n", data)
			s.WS.WriteMessage(websocket.TextMessage, data)
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

	Manager.Sessions = append(Manager.Sessions, sesh)

	return sesh.UUID.String()
}

func SessionExists(uuidString string) bool {
	for _, s := range Manager.Sessions {
		if s.UUID.String() == uuidString {
			return true
		}
	}
	return false
}

func AddWS(uuid string, c *websocket.Conn) error {

	for _, s := range Manager.Sessions {
		if s.UUID.String() == uuid {
			s.AddWS(c)
			return nil
		}
	}
	return fmt.Errorf("adding websocket to broker failed, uuid %s did not exist", uuid)
}
