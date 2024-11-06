package testnats_test

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
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				natsClient.Subscribe("test_subscribe", func(m *nats.Msg) {
					fmt.Println("NATS CLIENT, received msg")
					fmt.Println("  Header:", m.Header)
					fmt.Println("  Data:", string(m.Data))

					err := m.Respond([]byte("test_subscribe 222222"))
					Expect(err).ShouldNot(HaveOccurred())

					wg.Done()
				})
			}()

			time.Sleep(3 * time.Second)

			replay := natsHook.NewInbox()

			wg.Add(1)
			natsHook.Subscribe("test_subscribe", func(msg *nats.Msg) {
				fmt.Println("NATS HOOK, received msg")
				fmt.Println("Inbox:", replay)
				fmt.Println("Message:", string(msg.Data))

				//wg.Done()
			})

			fmt.Println("replay =", replay)
			err := natsHook.PublishRequest("test_subscribe", replay, []byte("test message 1111111"))
			Expect(err).ShouldNot(HaveOccurred())

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
