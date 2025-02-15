package main

import "encoding/json"

type BaseJson struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type BaseObject struct {
	Type string
	Data interface{}
}

type SubscriptionObject struct {
	Id         string
	TopicName  string
	BufferSize int
}

type UnsubscriptionObject struct {
	Id        string
	TopicName string
}

type PublishObject struct {
	TopicName string
	Message   string
}
