package main

import (
	"sync"
)

type Topic struct {
	name string

	sbMu        sync.RWMutex
	subscribers map[string]*Subscriber
}

func MakeTopic(name string) *Topic {
	return &Topic{
		name:        name,
		sbMu:        sync.RWMutex{},
		subscribers: make(map[string]*Subscriber),
	}
}

func (t *Topic) Publish(msg string) {
	t.sbMu.RLock()
	defer t.sbMu.RUnlock()

	// look if this needs go routines
	for _, subscriber := range t.subscribers {
		subscriber.ch <- msg

	}
}

func (t *Topic) AddSubscriber(sub *Subscriber) {
	t.sbMu.Lock()
	defer t.sbMu.Unlock()

	t.subscribers[sub.id] = sub
}

func (t *Topic) RemoveSubscriber(id string) {
	t.sbMu.Lock()
	defer t.sbMu.Unlock()
	delete(t.subscribers, id)
}

func (t *Topic) DoesSubscriberExist(subId string) bool {
	t.sbMu.RLock()
	defer t.sbMu.RUnlock()

	_, ok := t.subscribers[subId]
	return ok
}
