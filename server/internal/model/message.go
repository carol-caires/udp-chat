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
	MessageTypeNewMessage MessageType = "NEW_MESSAGE"
	MessageTypeDeleteMessage MessageType = "DELETE_MESSAGE"
	MessageTypeNewClient    MessageType = "NEW_CLIENT"
	MessageTypeDeleteClient MessageType = "DELETE_CLIENT"
)

var validMessageTypes = []MessageType{MessageTypeNewMessage, MessageTypeDeleteMessage, MessageTypeNewClient, MessageTypeDeleteClient}

func NewMessage(body string) (Message, string, error) {
	var message Message
	err := json.Unmarshal([]byte(body), &message)
	if err != nil {
		return Message{}, "", err
	}

	err = validateMessageType(message)
	return message, body, err
}

func validateMessageType(message Message) error {
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
