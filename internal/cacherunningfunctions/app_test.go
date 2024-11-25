package cacherunningfunctions_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCacheRunningFunction(t *testing.T) {

	cache, err := cacherunningfunctions.CreateCache(context.Background(), 10000)
	assert.NoError(t, err)

	var testStr string

	newId := uuid.New().String()
	cache.SetMethod(newId, func() bool {
		testStr = "test_string"

		fmt.Println("add method is started")

		return true
	})

	time.Sleep(time.Second * 7)
	assert.Equal(t, testStr, "test_string")
}
