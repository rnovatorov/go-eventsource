package eventsource

import "context"

type metadataContextKey struct{}

type Metadata map[string]interface{}

func WithMetadata(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, md)
}

func MetadataFromContext(ctx context.Context) Metadata {
	md, _ := ctx.Value(metadataContextKey{}).(Metadata)
	return md
}
