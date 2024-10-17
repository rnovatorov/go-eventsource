package sagacoordinator

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type RequestSender interface {
	SendRequest(context.Context, proto.Message) error
}
