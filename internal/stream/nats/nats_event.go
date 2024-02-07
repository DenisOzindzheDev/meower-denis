package nats

import (
	"bytes"
	"encoding/gob"
	"log"
	"meower-denis/internal/models"
	"meower-denis/internal/stream"

	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	nc                      *nats.Conn
	meowCreatedSubscription *nats.Subscription
	meowCreatedChan         chan stream.MeowCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}
func (e *NatsEventStore) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.meowCreatedSubscription != nil {
		e.meowCreatedSubscription.Unsubscribe()
	}
	close(e.meowCreatedChan)
}

func (e *NatsEventStore) PublishMeowCreated(meow models.Meow) error {
	m := stream.MeowCreatedMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return e.nc.Publish(m.Key(), data)
}
func (e *NatsEventStore) writeMessage(m stream.Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return b.Bytes(), nil
}
func (e *NatsEventStore) OnMeowCreated(f func(stream.MeowCreatedMessage)) (err error) {
	m := stream.MeowCreatedMessage{}
	e.meowCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		e.readMessage(msg.Data, &m)
		f(m)
	})
	return
}
func (mq *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (e *NatsEventStore) SubscribeMeowCreated() (<-chan stream.MeowCreatedMessage, error) {
	m := stream.MeowCreatedMessage{}
	e.meowCreatedChan = make(chan stream.MeowCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	e.meowCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan stream.MeowCreatedMessage)(e.meowCreatedChan), nil
}
