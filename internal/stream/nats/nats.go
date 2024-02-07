package nats

import (
	"meower-denis/internal/models"
	"meower-denis/internal/stream"
)

type EventStore interface {
	Close()
	PublishMeowCreated(meow models.Meow) error
	SubscribeMeowCreated() (<-chan stream.MeowCreatedMessage, error) //way 1 return channel
	OnMeowCreated(f func(stream.MeowCreatedMessage)) error           //way 2 return function
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishMeowCreated(meow models.Meow) error {
	return impl.PublishMeowCreated(meow)
}
func SubscribeMeowCreated() (<-chan stream.MeowCreatedMessage, error) {
	return impl.SubscribeMeowCreated()
}
func OnMeowCreated(f func(stream.MeowCreatedMessage)) error {
	return impl.OnMeowCreated(f)
}
