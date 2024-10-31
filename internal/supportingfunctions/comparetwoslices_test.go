package supportingfunctions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func TestCompareTwoSlice(t *testing.T) {
	mainList := []string{"green", "yellow", "red", "black", "white", "pink"}
	compareList := []string{"red", "white", "blue", "orange", "green"}

	result := supportingfunctions.CompareTwoSlices(mainList, compareList)
	fmt.Println("1. result:", result)
	assert.Equal(t, len(result), 2)

	result = supportingfunctions.CompareTwoSlices(compareList, mainList)
	fmt.Println("2. result:", result)
	assert.Equal(t, len(result), 3)

	result = supportingfunctions.CompareTwoSlices(mainList, []string{})
	fmt.Println("3. result:", result)
	assert.Equal(t, len(result), 0)

	result = supportingfunctions.CompareTwoSlices([]string{}, compareList)
	fmt.Println("4. result:", result)
	assert.Equal(t, len(result), 5)
}
