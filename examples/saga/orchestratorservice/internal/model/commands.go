package model

import (
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type BeginSaga struct {
	TaskDefinitions []*sagapb.TaskDefinition
}

type EndTask struct {
	ID     string
	Result *structpb.Value
}

type AbortTask struct {
	ID     string
	Reason *structpb.Value
}

type EndCompensation struct {
	ID string
}
