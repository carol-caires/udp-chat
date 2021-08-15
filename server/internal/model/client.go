package model

import "encoding/json"

type Client struct {
	Name string `json:"name"`
	IpAddress string `json:"address"`
}

func ParseClient(clientStr string) (Client, error) {
	var client Client
	err := json.Unmarshal([]byte(clientStr), &client)
	return client, err
}

func ParseClientsArray(clientsStr string) ([]Client, error) {
	var clients []Client
	err := json.Unmarshal([]byte(clientsStr), &clients)
	return clients, err
}
