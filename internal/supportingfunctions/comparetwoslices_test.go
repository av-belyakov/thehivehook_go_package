package supportingfunctions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func TestJoinTwoSlicesUniqValues(t *testing.T) {
	mainList := []string{"green", "yellow", "red", "black", "white", "pink"}
	compareList := []string{"red", "white", "blue", "orange", "green"}

	result := supportingfunctions.JoinTwoSlicesUniqValues(mainList, compareList)
	fmt.Println("1. result:", result)
	assert.Equal(t, len(result), 8)

	result = supportingfunctions.JoinTwoSlicesUniqValues(compareList, mainList)
	fmt.Println("2. result:", result)
	assert.Equal(t, len(result), 8)

	result = supportingfunctions.JoinTwoSlicesUniqValues(mainList, []string{})
	fmt.Println("3. result:", result)
	assert.Equal(t, len(result), 6)

	result = supportingfunctions.JoinTwoSlicesUniqValues([]string{}, compareList)
	fmt.Println("4. result:", result)
	assert.Equal(t, len(result), 5)

	result = supportingfunctions.JoinTwoSlicesUniqValues(
		[]string{"Webhook: send=\"Elasticsearch\"", "Webhook: send=\"MISP\"", "Webhook: send=\"ElasticsearchDB\""},
		[]string{"Webhook: send=\"MISP\""},
	)
	fmt.Println("5. result:", result)
	assert.Equal(t, len(result), 3)

	result = supportingfunctions.JoinTwoSlicesUniqValues(
		[]string{"Webhook: send=\"Elasticsearch\"", "Webhook: send=\"ElasticsearchDB\""},
		[]string{"Webhook: send=\"MISP\""},
	)
	fmt.Println("6. result:", result)
	assert.Equal(t, len(result), 3)

	result = supportingfunctions.JoinTwoSlicesUniqValues(
		[]string{
			"Webhook: send=\"Elasticsearch\"",
			"Sensor:id=\"8015632\"",
			"Webhook: send=\"ElasticsearchDB\"",
			"TheHivehook: send=\"MongoDB\"",
		},
		[]string{"Webhook: send=\"MISP\""},
	)
	fmt.Println("7. result:", result)
	assert.Equal(t, len(result), 5)

	result = supportingfunctions.JoinTwoSlicesUniqValues(
		[]string{
			"Webhook: send=\"Elasticsearch\"",
			"Sensor:id=\"8015632\"",
			"Webhook: send=\"MISP\"",
			"Webhook: send=\"ElasticsearchDB\"",
			"TheHivehook: send=\"MongoDB\"",
		},
		[]string{"Webhook: send=\"MISP\""},
	)
	fmt.Println("8. result:", result)
	assert.Equal(t, len(result), 5)

	result = supportingfunctions.JoinTwoSlicesUniqValues(
		[]string{},
		[]string{"Webhook: send=\"MISP\""},
	)
	fmt.Println("9. result:", result)
	assert.Equal(t, len(result), 1)
}
