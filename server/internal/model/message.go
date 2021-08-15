package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id string `json:"id"`
	MessageType MessageType `json:"type"`
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	From Client `json:"from"`
}

type MessageType string
const (
	MessageTypeMessage MessageType = "MESSAGE"
	MessageTypeNewClient MessageType = "NEW_CLIENT"
	MessageTypeCloseClient MessageType = "CLOSE_CLIENT"
)

func NewMessage(body string) (Message, string, error) {
	message := Message{Id: uuid.New().String()}
	err := json.Unmarshal([]byte(body), &message)
	return message, body, err
}
