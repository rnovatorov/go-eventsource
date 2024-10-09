package sagamodel

import (
	"google.golang.org/protobuf/types/known/structpb"
)

type BeginSaga struct {
	TaskDefinitions []*TaskDefinition
}

type EndTask struct {
	ID     string
	Result *structpb.Struct
}

type AbortTask struct {
	ID     string
	Reason *structpb.Struct
}

type EndCompensation struct {
	ID string
}
