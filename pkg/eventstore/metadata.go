package eventstore

import "context"

type metadataContextKey struct{}

type Metadata map[string]interface{}

func (m Metadata) CausationID() string {
	v, ok := m[CausationID]
	if !ok {
		return ""
	}
	causationID, _ := v.(string)
	return causationID
}

func WithMetadata(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, md)
}

func MetadataFromContext(ctx context.Context) Metadata {
	md, _ := ctx.Value(metadataContextKey{}).(Metadata)
	return md
}

const (
	CausationID = "X-Causation-ID"
)
