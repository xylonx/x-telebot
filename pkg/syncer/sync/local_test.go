package syncer_test

import (
	"os"
	"testing"
)

func TestDir(t *testing.T) {
	err := os.MkdirAll("configs/config", 0755)
	if err != nil {
		panic(err)
	}
}
