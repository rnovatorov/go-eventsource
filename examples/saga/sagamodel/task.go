package sagamodel

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type task struct {
	definition        *sagapb.TaskDefinition
	begun             bool
	ended             bool
	result            *structpb.Value
	aborted           bool
	abortReason       *structpb.Value
	compensationBegun bool
	compensationEnded bool
}

func (t *task) inProgress() bool {
	return t.begun && !t.aborted && !t.ended && !t.compensationBegun
}

func (t *task) compensationInProgress() bool {
	return t.compensationBegun && !t.compensationEnded
}

func (t *task) copy() *task {
	return &task{
		definition:        proto.Clone(t.definition).(*sagapb.TaskDefinition),
		begun:             t.begun,
		ended:             t.ended,
		result:            proto.Clone(t.result).(*structpb.Value),
		aborted:           t.aborted,
		abortReason:       proto.Clone(t.abortReason).(*structpb.Value),
		compensationBegun: t.compensationBegun,
		compensationEnded: t.compensationEnded,
	}
}
