package models

type Event struct {
	EventType string `json:"eventType"`
	Data      string `json:"data"`
}
