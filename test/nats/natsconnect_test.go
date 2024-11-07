package nats_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNatsConnect(t *testing.T) {
	prefix := "test"

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", "nats.cloud.gcm", 4222),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))

	assert.NoError(t, err)

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		if err != nil {
			assert.NoError(t, err)
		}

		fmt.Printf("module: '%s' the connection with NATS has been disconnected\n", prefix)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		if err != nil {
			assert.NoError(t, err)
		}

		fmt.Printf("module: '%s' the connection to NATS has been re-established\n", prefix)
	})

	nc.Subscribe("test_subscribe", func(m *nats.Msg) {
		fmt.Println("NATS CLIENT, received msg")
		fmt.Println("  Header:", m.Header)
		fmt.Println("  Data:", string(m.Data))

		m.Respond([]byte("test_subscribe 222222"))
	})

	time.Sleep(15 * time.Second)

	assert.Equal(t, true, true)
}
