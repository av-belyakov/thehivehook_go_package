package testthehiverequest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var _ = Describe("Testthehiverequest", Ordered, func() {
	var (
		rootDir string = "thehivehook_go_package"

		conf           *confighandler.ConfigApp
		chanTheHiveAPI chan<- thehiveapi.ReguestChannelTheHive

		chanDone chan struct{}

		errConf, errTheHiveApi error
	)

	BeforeAll(func() {
		chanDone = make(chan struct{})

		conf, errConf = confighandler.NewConfig(rootDir)
		confTheHive := conf.GetApplicationTheHive()

		//перед запуском теста установите переменную окружения GO_HIVEHOOK_THAPIKEY
		//с ключем-идентификатором, необходимым для авторизации в API TheHive,
		//командой export GO_HIVEHOOK_THAPIKEY=<api_key>

		logging := logginghandler.New()

		go func() {
			for {
				select {
				case <-chanDone:
					return

				case msg := <-logging.GetChan():
					fmt.Println("Log:", msg)
				}
			}
		}()

		chanTheHiveAPI, errTheHiveApi = thehiveapi.New(context.Background(), os.Getenv("GO_HIVEHOOK_THAPIKEY"), confTheHive.Host, confTheHive.Port, logging)
	})

	Context("Тест 0. Чтение конфигурационного файла", func() {
		It("При чтении конфигурационного файла не должно быть ошибок", func() {
			Expect(errConf).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1. Инициализация модуля взаимодействия с API TheHive", func() {
		It("При инициализации модуля не должно быть ошибок", func() {
			Expect(errTheHiveApi).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Выполнение запросов к TheHive", func() {
		It("Запрос на получения Observable должен быть успешно выполнен", func() {
			var (
				statusCode int
				rootId     string = "~86676517008" //caseId:35144
				myUuid     string = uuid.New().String()
				myUuidRes  string
				wg         sync.WaitGroup

				chanResObservable chan thehiveapi.ResponseChannelTheHive = make(chan thehiveapi.ResponseChannelTheHive)
				chanResTTL        chan thehiveapi.ResponseChannelTheHive = make(chan thehiveapi.ResponseChannelTheHive)
			)

			wg.Add(2)

			go func() {
				for res := range chanResObservable {
					myUuidRes = res.RequestId
					statusCode = res.StatusCode

					fmt.Println("--------- Observable ----------")
					fmt.Println("Resived Response")
					fmt.Println("RequestId:", res.RequestId)

					msg := []interface{}{}
					err := json.Unmarshal(res.Data, &msg)
					fmt.Println("ERROR:", err)
					fmt.Println("DATA:", msg)
				}

				wg.Done()
			}()
			go func() {
				for res := range chanResTTL {
					myUuidRes = res.RequestId
					statusCode = res.StatusCode

					fmt.Println("--------- TTL ----------")
					fmt.Println("Resived Response")
					fmt.Println("RequestId:", res.RequestId)

					msg := []interface{}{}
					err := json.Unmarshal(res.Data, &msg)
					fmt.Println("ERROR:", err)
					fmt.Println("DATA:", msg)
				}

				wg.Done()
			}()

			chanTheHiveAPI <- thehiveapi.ReguestChannelTheHive{
				RequestId:  myUuid,
				RootId:     rootId,
				Command:    "get_observables",
				ChanOutput: chanResObservable,
			}

			chanTheHiveAPI <- thehiveapi.ReguestChannelTheHive{
				RequestId:  myUuid,
				RootId:     rootId,
				Command:    "get_ttp",
				ChanOutput: chanResTTL,
			}

			wg.Wait()

			chanDone <- struct{}{}

			Expect(statusCode).Should(Equal(200))
			Expect(myUuidRes).Should(Equal(myUuid))
		})
		/*It("Запрос на получения TTP должен быть успешно выполнен", func() {

		})*/
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
