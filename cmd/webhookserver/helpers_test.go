package webhookserver_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/stretchr/testify/assert"
)

func TestGetter(t *testing.T) {
	jsonObj := fmt.Append(nil, `{
	"operation":  "update",
	"objectType": "case",
	"rootId":     "7yw9243yfw9g27rw",
	"object": {
		"createdAt": "2454642554",
		"caseId":    1234,
		"tags":      ["one_tag", "two_tag", "three_tag"]
	}}`)

	list := map[string]any{}

	t.Run("Test json unmarshal", func(t *testing.T) {
		err := json.Unmarshal(jsonObj, &list)
		assert.NoError(t, err)
	})

	t.Run("Test get objectType", func(t *testing.T) {
		value, err := webhookserver.GetObjectType(list)
		assert.NoError(t, err)
		assert.Equal(t, value, "case")
	})

	t.Run("Test get rootId", func(t *testing.T) {
		value, err := webhookserver.GetRootId(list)
		assert.NoError(t, err)
		assert.Equal(t, value, "7yw9243yfw9g27rw")
	})

	t.Run("Test get operation", func(t *testing.T) {
		value, err := webhookserver.GetOperation(list)
		assert.NoError(t, err)
		assert.Equal(t, value, "update")
	})

	t.Run("Test get caseId", func(t *testing.T) {
		caseId, err := webhookserver.GetCaseId(list)
		assert.NoError(t, err)

		assert.Equal(t, caseId, 1234)
	})
}
