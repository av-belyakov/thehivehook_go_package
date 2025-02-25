package main

//listenerbehindnats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

var (
	nc *nats.Conn
	f  *os.File

	chDone chan struct{} = make(chan struct{})

	err error
)

func ClientNATS(host string, port int) (*nats.Conn, error) {
	nc, err = nats.Connect(
		fmt.Sprintf("%s:%d", host, port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return nil, err
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println(err)
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println(err)
	})

	return nc, nil
}

func init() {
	nc, err = ClientNATS("nats.cloud.gcm", 4222)
	if err != nil {
		log.Panicln(err)
	}

	f, err = os.OpenFile(filepath.Join("internal", "listenerbehindnats", "case_test.log"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	go func() {
		log.Printf("system call:%+v", <-ctx.Done())

		nc.Close()
		f.Close()
		stop()
	}()

	nc.Subscribe("object.casetype", func(msg *nats.Msg) {
		ee := datamodels.CaseEventElement{}
		err = json.Unmarshal(msg.Data, &ee)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Received object type:'%s', root id:'%s' case id:'%d'\n", ee.ObjectType, ee.RootId, ee.Object.CaseId)

		str, err := supportingfunctions.NewReadReflectJSONSprint(msg.Data)
		if err != nil {
			log.Panicln(err)
		}

		_, err = f.WriteString(str)
		if err != nil {
			log.Panicln(err)
		}
	})

	fmt.Println("Start package listener NATS messages")

	<-chDone

	fmt.Println("Stop package")
}
