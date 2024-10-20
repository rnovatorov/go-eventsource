package orchestratorservice

import (
	"context"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type RequestSender interface {
	SendTaskRequest(context.Context, *sagapb.TaskRequest) error
	SendCompensationRequest(context.Context, *sagapb.CompensationRequest) error
}
