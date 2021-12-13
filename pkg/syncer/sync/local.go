package syncer

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	ErrorIsNotDirectory = errors.New("the input dir is not a directory")
	ErrorEmptyDirectory = errors.New("the dir is empty")
)

type LocalSyncer struct {
	directory string
}

var _ Synchronizer = &LocalSyncer{}

func NewLocalSyncer(dir string) (Synchronizer, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		// if directory not exists, create it recursively
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	if !stat.IsDir() {
		return nil, ErrorIsNotDirectory
	}

	return &LocalSyncer{directory: dir}, nil
}

func (s *LocalSyncer) Persistent(ctx context.Context, key string, data io.Reader) (string, error) {
	tmpLocation := path.Join(s.directory, key+".tmp")
	defer os.Remove(tmpLocation)

	location := path.Join(s.directory, key)

	file, err := os.Create(location)
	if err != nil {
		return "", err
	}

	// using default 32k buffer
	if _, err := io.Copy(file, data); err != nil {
		return "", err
	}

	if err := os.Rename(tmpLocation, location); err != nil {
		return "", err
	}

	return location, nil
}

func (s *LocalSyncer) PickOne(ctx context.Context) (string, error) {
	dir, err := os.Open(s.directory)
	if err != nil {
		return "", err
	}
	entries, err := dir.ReadDir(-1)
	if err != nil {
		return "", err
	}

	if len(entries) <= 0 {
		return "", ErrorEmptyDirectory
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	return path.Join(s.directory, entries[r.Intn(len(entries))].Name()), nil
}
