package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorageAcceptedCommands(t *testing.T) {
	storage, err := storage.NewStorageAcceptedCommands(
		storage.WithMaxSize(16),
		storage.WithMaxTtl(5),
		storage.WithTimeTick(2))
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	storage.Start(ctx)

	err = storage.SetObject("key_one", []byte(""))
	assert.NoError(t, err)

	_, ok := storage.GetObject("key_one")
	assert.True(t, ok)

	storage.SetObject("key_two", []byte(""))
	assert.NoError(t, err)
	storage.SetObject("key_three", []byte(""))
	assert.NoError(t, err)

	assert.Equal(t, len(storage.GetObjects()), 3)

	time.Sleep(7 * time.Second)

	assert.Equal(t, len(storage.GetObjects()), 0)

	cancel()
}
