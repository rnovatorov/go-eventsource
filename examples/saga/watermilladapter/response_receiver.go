package watermilladapter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/rnovatorov/go-eventsource/examples/saga/sagapb"
)

type UnmarshalFunc func([]byte, proto.Message) error

type ResponseHandler interface {
	HandlePerformTaskResponse(
		context.Context, *sagapb.PerformTaskResponse,
	) error
	HandlePerformCompensationResponse(
		context.Context, *sagapb.PerformCompensationResponse,
	) error
}

type ResponseReceiver struct {
	router        *message.Router
	handler       ResponseHandler
	unmarshalFunc UnmarshalFunc
	stopped       chan struct{}
}

type ResponseReceiverParams struct {
	RouterConfig   message.RouterConfig
	Logger         *slog.Logger
	Subscriber     message.Subscriber
	SubscribeTopic string
	Handler        ResponseHandler
	UnmarshalFunc  UnmarshalFunc
}

func StartResponseReceiver(
	ctx context.Context, p ResponseReceiverParams,
) (*ResponseReceiver, error) {
	router, err := message.NewRouter(
		p.RouterConfig, watermill.NewSlogLogger(p.Logger))
	if err != nil {
		return nil, fmt.Errorf("new message router: %w", err)
	}

	if p.UnmarshalFunc == nil {
		p.UnmarshalFunc = protojson.Unmarshal
	}

	r := &ResponseReceiver{
		router:        router,
		handler:       p.Handler,
		unmarshalFunc: p.UnmarshalFunc,
		stopped:       make(chan struct{}),
	}

	router.AddMiddleware(middleware.Retry{
		Logger: watermill.NewSlogLogger(p.Logger),
	}.Middleware)

	router.AddNoPublisherHandler("response_receiver", p.SubscribeTopic,
		p.Subscriber, r.processMessage)

	go func() {
		defer close(r.stopped)
		_ = router.Run(ctx)
	}()

	return r, nil
}

func (r *ResponseReceiver) Stop() {
	_ = r.router.Close()
	<-r.stopped
}

func (r *ResponseReceiver) processMessage(msg *message.Message) error {
	var payload anypb.Any
	if err := r.unmarshalFunc(msg.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal any: %w", err)
	}

	response, err := payload.UnmarshalNew()
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	switch resp := response.(type) {
	case *sagapb.PerformTaskResponse:
		if err := r.handler.HandlePerformTaskResponse(
			msg.Context(), resp,
		); err != nil {
			return fmt.Errorf("handle %T: %w", resp, err)
		}
		return nil
	case *sagapb.PerformCompensationResponse:
		if err := r.handler.HandlePerformCompensationResponse(
			msg.Context(), resp,
		); err != nil {
			return fmt.Errorf("handle %T: %w", resp, err)
		}
		return nil
	default:
		return fmt.Errorf("%w: %T", errUnexpectedResponseType, resp)
	}
}
