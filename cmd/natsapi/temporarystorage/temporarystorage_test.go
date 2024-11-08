package temporarystoarge_test

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	temporarystoarge "github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/temporarystorage"
)

func TestTemporaryStorage(t *testing.T) {
	ts, err := temporarystoarge.NewTemporaryStorage(context.Background(), 7)
	assert.NoError(t, err)

	err = ts.SetService("34d11", "TestService")
	assert.Error(t, err)

	id := ts.NewCell()

	//**** Service
	err = ts.SetService(id, "TestService")
	assert.NoError(t, err)

	value, err := ts.GetService(id)
	assert.NoError(t, err)
	assert.Equal(t, value, "TestService")

	err = ts.SetService(id, "_testService_")
	assert.NoError(t, err)

	value, err = ts.GetService(id)
	assert.NoError(t, err)
	assert.Equal(t, value, "_testService_")

	//**** RootId
	err = ts.SetRootId(id, "~6438882")
	assert.NoError(t, err)

	value, err = ts.GetRootId(id)
	assert.NoError(t, err)
	assert.Equal(t, value, "~6438882")

	//**** CaseId
	err = ts.SetCaseId(id, "1353543")
	assert.NoError(t, err)

	value, err = ts.GetCaseId(id)
	assert.NoError(t, err)
	assert.Equal(t, value, "1353543")

	//**** NsMsg
	err = ts.SetNsMsg(id, &nats.Msg{})
	assert.NoError(t, err)

	_, err = ts.GetNsMsg(id)
	assert.NoError(t, err)

	//******************** another id ********************
	anotherId := ts.NewCell()

	//**** Service
	err = ts.SetService(anotherId, "TestAnyIdService")
	assert.NoError(t, err)

	value, err = ts.GetService(anotherId)
	assert.NoError(t, err)
	assert.Equal(t, value, "TestAnyIdService")

	//**** NsMsg
	err = ts.SetNsMsg(anotherId, &nats.Msg{})
	assert.NoError(t, err)

	_, err = ts.GetNsMsg(anotherId)
	assert.NoError(t, err)

	ts.DeleteElement(anotherId)
	_, err = ts.GetCaseId(anotherId)
	assert.Error(t, err)

	//********* finaly **********
	value, err = ts.GetRootId(id)
	assert.NoError(t, err)
	assert.Equal(t, value, "~6438882")
}
