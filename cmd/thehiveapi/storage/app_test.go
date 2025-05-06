package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/storage"
	"github.com/stretchr/testify/assert"
)

func TestStorageFoundObjects(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	sfo, err := storage.NewStorageFoundObjects(
		storage.WithMaxSize(3),
		storage.WithMaxTtl(10),
		storage.WithTimeTick(2))
	assert.NoError(t, err)

	sfo.Start(ctx)

	err = sfo.SetObject("one", []byte("one element"))
	assert.NoError(t, err)

	_, ok := sfo.GetObject("one")
	assert.True(t, ok)

	err = sfo.SetObject("two", []byte("two element"))
	assert.NoError(t, err)
	assert.Equal(t, len(sfo.GetObjects()), 2)

	err = sfo.SetObject("three", []byte("three element"))
	assert.NoError(t, err)
	assert.Equal(t, len(sfo.GetObjects()), 3)

	time.Sleep(1 * time.Second)

	_, ok = sfo.GetObject("three")
	assert.True(t, ok)

	err = sfo.SetObject("four", []byte("four element"))
	assert.NoError(t, err)
	assert.Equal(t, len(sfo.GetObjects()), 3)

	time.Sleep(12 * time.Second)

	list := sfo.GetObjects()
	t.Log("List objects:", list)

	assert.Equal(t, len(list), 0)

	cancel()
}
