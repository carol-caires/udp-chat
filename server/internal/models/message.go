package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id string `json:"id"`
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	From Client `json:"from"`
}

func NewMessage(client Client, body string) (Message, string) {
	messageId := uuid.New()
	message := Message{
		Id: messageId.String(),
		Body:      body,
		CreatedAt: time.Now(),
		From:      client,
	}
	messageBytes, _ := json.Marshal(message)
	return message, string(messageBytes)
}
