package models

type MessageType int

const (
	MessageTypeUpdate MessageType = iota
	MessageTypeDelete
)

// ProductUpdate message identifies items to index
type Message struct {
	Type     MessageType `json:"messageType"`
	Products []Product   `json:"products"`
}
