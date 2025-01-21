package supportingfunctions_test

import (
	"fmt"
	"runtime"
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

	t.Log(One())

	assert.Equal(t, len(commonList), 9)
}

func One() string {
	_, f, l, _ := runtime.Caller(1)

	fmt.Println(Two())

	return fmt.Sprintf("function 'One', file:%s, line:%d", f, l)
}

func Two() string {
	_, f, l, _ := runtime.Caller(1)

	return fmt.Sprintf("function 'Two', file:%s, line:%d", f, l)
}
