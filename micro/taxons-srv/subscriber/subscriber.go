package subscriber

import (
	"context"
	"github.com/micro/go-log"

	proto "btdxcx.com/micro/taxons-srv/proto/imp"
)

// Receiver receive Struct as Subscriber
type Receiver struct{}

// Handle receive Struct as Subscriber
func (r *Receiver) Handle(ctx context.Context, msg *proto.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

// Handler receive Function as Subscriber
func Handler(ctx context.Context, msg *proto.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
