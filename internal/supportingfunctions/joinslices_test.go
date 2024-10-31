package supportingfunctions_test

import (
	"fmt"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"github.com/stretchr/testify/assert"
)

func TestJoinSlices(t *testing.T) {
	list1 := []string{"green", "yellow", "red", "black", "white"}
	list2 := []string{"red", "blue", "orange", "green"}
	list3 := []string{"yellow", "pink", "grey"}

	commonList := supportingfunctions.JoinSlises[string](list1, list2, list3)
	fmt.Println("common list:", commonList)

	assert.Equal(t, len(commonList), 9)
}
