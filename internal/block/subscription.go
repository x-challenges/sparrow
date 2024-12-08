package block

import "github.com/x-challenges/raven/broadcaster"

// Subscription
type Subscription = broadcaster.Broadcaster[*Block]

// Listener
type Listener = broadcaster.Listener[*Block]

// NewSubscription
func NewSubscription() Subscription {
	return broadcaster.New[*Block]()
}
