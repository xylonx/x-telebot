package syncer

import (
	"context"
	"io"
)

// Synchronizer - sync remote data into local
type Synchronizer interface {
	// Persistent - save data
	Persistent(ctx context.Context, key string, data io.Reader) (location string, err error)

	// PickOne - randomly pick one
	PickOne(ctx context.Context) (location string, err error)
}
