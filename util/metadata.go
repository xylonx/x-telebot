package util

import "context"

type Metadata struct{}

func ContextWithMetadata(ctx context.Context, data map[string]string) context.Context {
	return context.WithValue(ctx, Metadata{}, data)
}

func MetadataFromContext(ctx context.Context) map[string]string {
	if d, ok := ctx.Value(Metadata{}).(map[string]string); ok {
		return d
	}
	return make(map[string]string)
}
