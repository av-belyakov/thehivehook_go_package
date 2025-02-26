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
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

var (
	nc     *nats.Conn
	fc, fa *os.File

	chDone chan struct{} = make(chan struct{})

	err error
)

type Element struct {
	Source string                      `json:"source"`
	Event  datamodels.CaseEventElement `json:"event"`
}

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

	fc, err = os.OpenFile(filepath.Join("internal", "listenerbehindnats", "case_test.log"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln(err)
	}

	fa, err = os.OpenFile(filepath.Join("internal", "listenerbehindnats", "alert_test.log"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)
	defer stop()

	go func() {
		<-ctx.Done()

		nc.Close()
		fc.Close()
		fa.Close()
	}()

	ee := Element{}

	//этот модуль может принимать несколько одинаковых сообщений от разных
	//источников, если например, запущены две копии thehivehook_go, одна для тестов
	//локально, а другая может быть развернута в докере
	//по этому стоит поменять наименование подписки, что бы она была только для
	//локального модуля

	//для кейсов
	nc.Subscribe("object.casetype.local", func(msg *nats.Msg) {
		err = json.Unmarshal(msg.Data, &ee)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Received Case object case id:'%d', root id:'%s'\n", ee.Event.Object.CaseId, ee.Event.RootId)

		str, err := supportingfunctions.NewReadReflectJSONSprint(msg.Data)
		if err != nil {
			log.Panicln(err)
		}

		_, err = fc.WriteString(fmt.Sprintf("---------------\n%s\n", str))
		if err != nil {
			log.Panicln(err)
		}
	})

	//для алертов
	nc.Subscribe("object.alerttype.local", func(msg *nats.Msg) {
		err = json.Unmarshal(msg.Data, &ee)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Received Alert object root id:'%s'\n", ee.Event.RootId)

		str, err := supportingfunctions.NewReadReflectJSONSprint(msg.Data)
		if err != nil {
			log.Panicln(err)
		}

		_, err = fa.WriteString(fmt.Sprintf("---------------\n%s\n", str))
		if err != nil {
			log.Panicln(err)
		}
	})

	fmt.Println("Start package listener NATS messages")

	<-chDone

	fmt.Println("Stop package")
}
