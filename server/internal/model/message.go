package model

import (
	"encoding/json"
	"errors"
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

var validMessageTypes = []MessageType{MessageTypeMessage, MessageTypeNewClient, MessageTypeCloseClient}

func NewMessage(body string) (*Message, string, error) {
	var message *Message
	err := json.Unmarshal([]byte(body), message)

	if message != nil {
		err = validateMessageType(message)
		if err != nil {
			return nil, "", err
		}
	}

	return message, body, err
}

func validateMessageType(message *Message) error {
	message.Id = uuid.New().String()

	var isTypeValid bool
	for _, t := range validMessageTypes {
		if t == message.MessageType {
			isTypeValid = true
		}
	}
	if !isTypeValid {
		return errors.New("impossible to identify message type")
	}
	return nil
}
