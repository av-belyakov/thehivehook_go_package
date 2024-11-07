package nats_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	NATS_HOST = "nats.cloud.gcm"
	NATS_PORT = 4222

	CACHETTL         = 360
	SENDER_CASE      = "object.casetype"
	SENDER_ALERT     = "object.alerttype"
	LISTENER_COMMAND = "object.commandstype"
)

var (
	natsHook, natsClient *nats.Conn

	errNatsHook, errNatsClient error
)

func CreateNatsConnect(prefix, host string, port int) (*nats.Conn, error) {
	var (
		nc  *nats.Conn
		err error
	)

	nc, err = nats.Connect(
		fmt.Sprintf("%s:%d", host, port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return nc, err
	}

	fmt.Println("func 'CreateNatsConnect', prefix:", prefix)

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		if err != nil {
			fmt.Println(err)
			fmt.Printf("module: '%s' the connection with NATS has been disconnected %s\n", prefix, err.Error())

			return
		}

		fmt.Printf("module: '%s' the connection with NATS has been disconnected\n", prefix)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		if err != nil {
			fmt.Printf("module: '%s' the connection to NATS has been re-established %s\n", prefix, err.Error())

			return
		}

		fmt.Printf("module: '%s' the connection to NATS has been re-established\n", prefix)
	})

	return nc, nil
}

var _ = Describe("Testnats", Ordered, func() {
	BeforeAll(func() {

		natsHook, errNatsHook = CreateNatsConnect("HOOK", NATS_HOST, NATS_PORT)
		natsClient, errNatsClient = CreateNatsConnect("CLIENT", NATS_HOST, NATS_PORT)
	})

	Context("Тест 1. Подключение к NATS", func() {
		It("При подключении к NATS natsHook не должно быть ошибки", func() {
			Expect(errNatsHook).ShouldNot(HaveOccurred())
		})

		It("При подключении к NATS natsClient не должно быть ошибки", func() {
			Expect(errNatsClient).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Отправка и прием сообщений", func() {
		It("Передаваемое от natsHook сообщение должно быть успешно получено", func() {
			var (
				wg sync.WaitGroup

				msgTestCase    string = "test message CASE"
				msgTestAlert   string = "test message ALERT"
				msgTestCommand string = "test message SEND TEST command"
			)
			subscribersSenders := [...]string{SENDER_CASE, SENDER_ALERT}

			wg.Add(len(subscribersSenders))
			for _, v := range subscribersSenders {
				go func(subscriber string) {
					natsClient.Subscribe(subscriber, func(m *nats.Msg) {
						msg := string(m.Data)
						fmt.Printf("NATS CLIENT, subscribe:'%s' received msg\n", v)
						fmt.Println("  Data:", msg)

						if subscriber == "object.casetype" {
							Expect(msg).Should(Equal(msgTestCase))
						} else if subscriber == "object.alerttype" {
							Expect(msg).Should(Equal(msgTestAlert))
						} else {
							Expect(false).Should(BeTrue())
						}

						wg.Done()
					})
				}(v)
			}

			wg.Add(1)
			go func() {
				natsHook.Subscribe(LISTENER_COMMAND, func(m *nats.Msg) {
					msg := string(m.Data)
					fmt.Printf("NATS HOOK, subscribe:'%s' received msg\n", LISTENER_COMMAND)
					fmt.Println("  Data:", msg)

					Expect(msg).Should(Equal(msgTestCommand))

					err := natsHook.Publish(m.Reply, []byte(fmt.Sprintf("the command '%s' was executed successfully", msg)))
					Expect(err).ShouldNot(HaveOccurred())

					wg.Done()
				})
			}()

			time.Sleep(3 * time.Second)

			//send case
			err := natsHook.Publish(SENDER_CASE, []byte(msgTestCase))
			Expect(err).ShouldNot(HaveOccurred())

			//send alert
			err = natsHook.Publish(SENDER_ALERT, []byte(msgTestAlert))
			Expect(err).ShouldNot(HaveOccurred())

			//send command
			msg, err := natsClient.Request(LISTENER_COMMAND, []byte(msgTestCommand), time.Second*3)
			//msg, err := nc.RequestWithContext(ctx, "foo", []byte("bar"))
			//err = natsClient.PublishRequest(LISTENER_COMMAND, natsHook.NewInbox(), []byte(msgTestCommand))
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println("---- RESPONSE:", string(msg.Data))
			// Synchronous subscriber with context
			//sub, err := nc.SubscribeSync("foo")
			//msg, err := sub.NextMsgWithContext(ctx)

			wg.Wait()

			Expect(true).Should(BeTrue())
		})
	})

	AfterAll(func() {
		natsHook.Close()
		natsClient.Close()
	})

	/*
		Context("", func(){
			It("", func ()  {

			})
		})
	*/
})
