package counterelements_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/counterelements"
)

func TestCounterElements(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ce := counterelements.New(3)
	ce.Start(ctx)

	ce.Add("key-1")
	ce.Add("key-2")
	ce.Add("key-1")
	ce.Add("key-1")
	assert.Equal(t, ce.Get("key-1"), 3)
	assert.Equal(t, ce.Get("key-2"), 1)

	ce.Done("key-1")
	assert.Equal(t, ce.Get("key-1"), 2)

	ce.Done("key-2")
	assert.Equal(t, ce.Get("key-2"), 0)
	ce.Done("key-2")
	assert.Equal(t, ce.Get("key-2"), 0)

	assert.Equal(t, ce.Get("key-3"), -1)
	ce.Done("key-3")
	assert.Equal(t, ce.Get("key-3"), -1)

	assert.Equal(t, ce.Size(), 2)

	time.Sleep(3 * time.Second)

	assert.Equal(t, ce.Size(), 0)

	ce.Add("key-3")
	ce.Add("key-4")
	ce.Add("key-5")
	assert.Equal(t, ce.Size(), 3)

	t.Cleanup(func() {
		cancel()
	})
}
