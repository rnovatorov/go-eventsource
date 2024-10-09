package sagamodel

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type task struct {
	arguments         *structpb.Struct
	dependencies      []string
	begun             bool
	ended             bool
	result            *structpb.Struct
	aborted           bool
	abortReason       *structpb.Struct
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
		arguments:         proto.Clone(t.arguments).(*structpb.Struct),
		dependencies:      t.dependencies,
		begun:             t.begun,
		ended:             t.ended,
		result:            proto.Clone(t.result).(*structpb.Struct),
		aborted:           t.aborted,
		abortReason:       proto.Clone(t.abortReason).(*structpb.Struct),
		compensationBegun: t.compensationBegun,
		compensationEnded: t.compensationEnded,
	}
}
