package natssendmessage__test

import (
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	nc, err := nats.Connect("192.168.9.208:4222",
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println(err)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println(err)
	})

	err = nc.Publish("object.casetype.local", []byte("ssddsdddvcbbngn message for check test"))
	assert.NoError(t, err)

	nc.Flush()

	err = nc.Drain()
	assert.NoError(t, err)

}

/*func TestManyConnect(t *testing.T) {
	for i := 0; i < 2024; i++ {
		nc, err := nats.Connect("192.168.9.208:4222",
			nats.MaxReconnects(-1),
			nats.ReconnectWait(3*time.Second))
		if err != nil {
			log.Fatalln(err)
		}

		if err = nc.Publish("object.casetype.local", []byte(fmt.Sprintf("message for check test, num connect: %d", i))); err != nil {
			fmt.Println(err)
		}
	}

	chS := make(chan os.Signal, 1)
	signal.Notify(chS, os.Interrupt, syscall.SIGINT)

	<-chS
}*/
