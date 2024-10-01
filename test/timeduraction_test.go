package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeDuraction(t *testing.T) {
	var ttl int = 5

	timeToLive, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	fmt.Println("time to live:", timeToLive)
	assert.Nil(t, err)
	assert.Equal(t, timeToLive, (5 * time.Second))
}
