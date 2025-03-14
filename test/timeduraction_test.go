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

	currentTime := time.Now()

	fmt.Println("Time:", int(currentTime.Month()))

	assert.Equal(t, currentTime.Year(), 2025)
}

func TestSub(t *testing.T) {
	timeStart := time.Date(2025, time.March, 13, 21, 32, 0, 0, time.Local)

	t.Log("count hours:", int(time.Since(timeStart).Hours()))
}
/*docker build -t gitlab.cloud.gcm:5050/a.belyakov/thehivehook_go_package:test \
      --build-arg VERSION=1.2.3 \
      --build-arg BRANCH=development \
      --build-arg STATUS=development .