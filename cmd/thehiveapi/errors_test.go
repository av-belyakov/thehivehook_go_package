package thehiveapi_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/stretchr/testify/assert"
)

func TestErrorAsType(t *testing.T) {
	errInfo := thehiveapi.ErrorInformation{"any information"}
	errNew := errors.New("new error for test")

	e, isExist := errors.AsType[thehiveapi.ErrorInformation](errInfo)
	fmt.Println("Error:", e)
	assert.True(t, isExist)

	e, isExist = errors.AsType[thehiveapi.ErrorInformation](errNew)
	fmt.Println("Error:", e)
	assert.False(t, isExist)
}
