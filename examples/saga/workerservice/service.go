package workerservice

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) HandleTaskRequest(
	ctx context.Context, req *sagapb.TaskRequest,
) (*sagapb.TaskResponse, error) {
	time.Sleep(3 * time.Second)

	if req.TaskId == "task3" {
		return &sagapb.TaskResponse{
			SagaId:          req.SagaId,
			TaskId:          req.TaskId,
			TaskAborted:     true,
			TaskResult:      nil,
			TaskAbortReason: structpb.NewStringValue(req.TaskId + " abort reason"),
		}, nil
	}

	return &sagapb.TaskResponse{
		SagaId:          req.SagaId,
		TaskId:          req.TaskId,
		TaskAborted:     false,
		TaskResult:      structpb.NewStringValue(req.TaskId + " result"),
		TaskAbortReason: nil,
	}, nil
}

func (s *Service) HandleCompensationRequest(
	ctx context.Context, req *sagapb.CompensationRequest,
) (*sagapb.CompensationResponse, error) {
	time.Sleep(3 * time.Second)

	return &sagapb.CompensationResponse{
		SagaId: req.SagaId,
		TaskId: req.TaskId,
	}, nil
}
