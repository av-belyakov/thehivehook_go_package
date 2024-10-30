package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeTicker(t *testing.T) {
	var num int

	for range time.Tick(1 * time.Second) {
		num++

		if num == 13 {
			break
		}

		if num < 7 {
			time.Sleep(3 * time.Second)
		}

		fmt.Println("num:", num)
	}

	assert.Equal(t, num, 13)
}
