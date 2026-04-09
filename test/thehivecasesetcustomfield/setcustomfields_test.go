package thehivecasesetcustomfield_test

import (
	"errors"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSetCustomField(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	theHiveApi, err := thehiveapi.NewTheHiveApi(
		thehiveapi.WithAPIKey_New(os.Getenv("GO_HIVEHOOK_THAPIKEY")),
		thehiveapi.WithHost_New("thehive.cloud.gcm"),
		thehiveapi.WithPort_New(9000),
		thehiveapi.WithNameRegionalObject_New("test-object"),
	)
	assert.NoError(t, err)

	command := thehiveapi.RequestCommand{
		Service:   "test-service",
		Command:   "set_case_custom_field",
		RootId:    "~93109592216",
		FieldName: "misp-event-id.string",
		Value:     "1222344566",
	}

	t.Run("Test 1. Поля не существует", func(t *testing.T) {
		_, statusCode, err := theHiveApi.AddCaseCustomFields(t.Context(), command)
		assert.Error(t, err)
		assert.Equal(t, statusCode, http.StatusNotFound)
	})

	t.Run("Test 2. Поле существует, но значение отличается от переданного в команде", func(t *testing.T) {
		command.RootId = "~90646261816"
		command.Value = "1222344566"

		_, statusCode, err := theHiveApi.AddCaseCustomFields(t.Context(), command)
		assert.NoError(t, err)
		assert.Equal(t, statusCode, http.StatusOK)

		//fmt.Println("Response:", string(b))
	})

	t.Run("Test 3. Поле существует, но значение соответствует значению в команде", func(t *testing.T) {
		command.RootId = "~90646261816"
		command.Value = "142243"

		_, statusCode, err := theHiveApi.AddCaseCustomFields(t.Context(), command)
		assert.Error(t, err)
		_, ok := errors.AsType[thehiveapi.ErrorInformation](err)
		assert.True(t, ok)
		assert.Equal(t, statusCode, http.StatusNotModified)

		//fmt.Println("Response:", string(b))
	})

	t.Cleanup(func() {
		os.Unsetenv("GO_HIVEHOOK_THAPIKEY")
	})
}
