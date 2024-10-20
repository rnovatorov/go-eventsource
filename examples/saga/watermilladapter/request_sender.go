package watermilladapter

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type MarshalFunc func(proto.Message) ([]byte, error)

type RequestSender struct {
	publisher message.Publisher
}

type RequestSenderParams struct {
	Publisher message.Publisher
}

func NewRequestSender(p RequestSenderParams) *RequestSender {
	return &RequestSender{
		publisher: p.Publisher,
	}
}

func (s *RequestSender) SendRequest(ctx context.Context, req proto.Message) error {
	payload, err := protojson.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	msg := message.NewMessage(uuid.NewString(), payload)
	msg.SetContext(ctx)

	// TODO
	topic := "test"

	if err := s.publisher.Publish(topic, msg); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}
