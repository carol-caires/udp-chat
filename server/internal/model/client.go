package model

import "encoding/json"

type Client struct {
	Name string `json:"name"`
	IpAddress string `json:"address"`
}

func ParseClientsArray(clientsStr string) ([]Client, error) {
	var clients []Client
	err := json.Unmarshal([]byte(clientsStr), &clients)
	return clients, err
}
