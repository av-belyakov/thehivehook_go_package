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

	var (
		chTest       chan string = make(chan string)
		chResultFunc chan bool   = make(chan bool)
	)

	idOne := uuid.New().String()
	cache.SetMethod(idOne, func(count int) bool {
		chTest <- "test_string"

		fmt.Println("add method is started, attempt number:", count)

		return true
	})

	idTwo := uuid.New().String()
	cache.SetMethod(idTwo, func(count int) bool {
		fmt.Println("function two, count attempts:", count)

		if count == 3 {
			fmt.Println("function running is:", count, " attempt")
			chResultFunc <- true

			return true
		}

		//chResultFunc <- false

		return false
	})

	time.Sleep(time.Second * 7)
	assert.Equal(t, <-chTest, "test_string")

	time.Sleep(time.Second * 17)
	assert.True(t, <-chResultFunc)
}
