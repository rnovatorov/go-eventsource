package sagaservice

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/rnovatorov/go-eventsource/pkg/saga/sagamodel"
)

type SagaContext interface {
	context.Context
	SagaID() string
	TaskArguments(string) proto.Message
	TaskResult(string) proto.Message
}

type sagaContext struct {
	context.Context
	saga *sagamodel.SagaAggregate
}

func (c *sagaContext) SagaID() string {
	return c.saga.ID()
}

func (c *sagaContext) TaskArguments(id string) proto.Message {
	return c.saga.Root().TaskArguments(id)
}

func (c *sagaContext) TaskResult(id string) proto.Message {
	return c.saga.Root().TaskResult(id)
}
