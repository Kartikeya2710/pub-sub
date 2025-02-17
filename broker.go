package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Broker struct {
	listenAddr string

	tpMu   sync.RWMutex
	topics map[string]*Topic
}

func NewBroker(listenAddr string) *Broker {
	return &Broker{
		listenAddr: listenAddr,
		tpMu:       sync.RWMutex{},
		topics:     make(map[string]*Topic),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (b *Broker) StartBroker() {
	http.HandleFunc("/ws", b.wsHandler)
	fmt.Println("Your broker is running on:", b.listenAddr)
	err := http.ListenAndServe(b.listenAddr, nil)
	if err != nil {
		fmt.Println("Error starting the broker. Exiting...")
		return
	}
}

func (b *Broker) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading http connection to websockets")
		return
	}

	defer conn.Close()
	var jsonMessage BaseJson

	for {
		err := conn.ReadJSON(&jsonMessage)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("IsUnexpectedCloseError()", err)
			}
			break
		}

		switch jsonMessage.Type {
		case "subscribe":
			var subObj SubscriptionObject
			json.Unmarshal(jsonMessage.Data, &subObj)
			b.subscribe(subObj, conn)
		case "unsubscribe":
			var unsubObj UnsubscriptionObject
			json.Unmarshal(jsonMessage.Data, &unsubObj)
			b.unsubscribe(unsubObj)
		case "publish":
			var pubObj PublishObject
			json.Unmarshal(jsonMessage.Data, &pubObj)
			b.publish(pubObj)
		default:
			fmt.Println("The received message does not fit into any of the categories defined")
		}
	}
}

func (b *Broker) subscribe(subObj SubscriptionObject, conn *websocket.Conn) {
	topic := b.getTopic(subObj.TopicName)
	if topic.DoesSubscriberExist(subObj.Id) {
		fmt.Printf("%s is already subscribed on the topic %s\n", subObj.Id, subObj.TopicName)
		return
	}

	sub := &Subscriber{
		id:         subObj.Id,
		conn:       conn,
		bufferSize: subObj.BufferSize,
		ch:         make(chan string, subObj.BufferSize),
	}
	topic.AddSubscriber(sub)
	go sub.StartSubscriber()
}

func (b *Broker) unsubscribe(unsubObj UnsubscriptionObject) {
	topic, ok := b.topicExists(unsubObj.TopicName)
	if !ok {
		fmt.Printf("No topic found with name %s. No unsubscription needed", unsubObj.TopicName)
		return
	}

	topic.RemoveSubscriber(unsubObj.Id)
	fmt.Printf("%s subscriber removed from topic %s", unsubObj.Id, topic.name)
}

func (b *Broker) publish(pubObj PublishObject) {
	topic := b.getTopic(pubObj.TopicName)
	topic.Publish(pubObj.Message)
}

// gets the topic if it exists otherwise creates one
func (b *Broker) getTopic(topicName string) *Topic {
	existingTopic, ok := b.topicExists(topicName)
	if ok {
		return existingTopic
	}

	topic := MakeTopic(topicName)
	b.topics[topicName] = topic

	return topic
}

// returns the topic and the boolean telling whether it exists in the broker
func (b *Broker) topicExists(topicName string) (*Topic, bool) {
	b.tpMu.Lock()
	defer b.tpMu.Unlock()

	topic, ok := b.topics[topicName]
	return topic, ok
}
