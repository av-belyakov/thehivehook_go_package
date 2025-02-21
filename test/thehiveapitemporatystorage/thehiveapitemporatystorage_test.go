package thehiveapitemporatystorage_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var (
	once        sync.Once
	cacheRunCom CacheRunningMethods

	//chanCacheMethodResponse chan CacheMethodResponse
)

type CacheMethodResponse struct {
	StatusCode   int
	Error        error
	BodyResponse []byte
}

type storageParameters struct {
	timeExpiry  time.Time
	cacheMethod func() bool
	//cacheMethod func(ctx context.Context, rootId string, i interface{}, chanRes chan<- CacheMethodResponse) bool
}

type cacheStorageParameters struct {
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

type CacheRunningMethods struct {
	ttl          time.Duration
	cacheStorage cacheStorageParameters
}

func New(ctx context.Context, ttl int) (*CacheRunningMethods, error) {
	cacheRunCom = CacheRunningMethods{}

	if ttl < 5 || ttl > 86400 {
		return &cacheRunCom, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	var err error
	once.Do(func() {
		timeToLive, newErr := time.ParseDuration(fmt.Sprintf("%ds", ttl))
		if newErr != nil {
			err = newErr

			return
		}

		cacheRunCom.ttl = timeToLive
		cacheRunCom.cacheStorage = cacheStorageParameters{
			storages: make(map[string]storageParameters),
		}

		go cacheRunCom.automaticExecutionMethods(ctx)
		//go checkLiveTime(&cacheRunCom)
	})

	return &cacheRunCom, err
}

// checkLiveTime удаляет устаревшую временную информацию
/*func checkLiveTime(crc *CacheRunningMethods) {
	for range time.Tick(5 * time.Second) {
		go func() {
			crc.cacheStorage.mutex.Lock()
			defer crc.cacheStorage.mutex.Unlock()

			for k, v := range crc.cacheStorage.storages {
				if v.timeExpiry.Before(time.Now()) {
					delete(crc.cacheStorage.storages, k)
				}
			}
		}()
	}
}*/

func (crm *CacheRunningMethods) automaticExecutionMethods(ctx context.Context) {
	tick := time.NewTicker(5 * time.Second)

	go func(ctx context.Context, tick *time.Ticker) {
		<-ctx.Done()
		tick.Stop()
	}(ctx, tick)

	for range tick.C {
		crm.cacheStorage.mutex.Lock()
		for k, v := range crm.cacheStorage.storages {
			//удаляем если записи слишком старые
			if v.timeExpiry.Before(time.Now()) {
				delete(crm.cacheStorage.storages, k)
			}

			if v.cacheMethod() {
				delete(crm.cacheStorage.storages, k)
			}
		}
		crm.cacheStorage.mutex.Unlock()
	}
}

// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningMethods) SetMethod(id string, f func() bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry:  time.Now().Add(crm.ttl),
		cacheMethod: f,
	}

	return id
}

// GetMethod возвращает данные по полученому id
func (crm *CacheRunningMethods) GetMethod(id string) (func() bool, bool) {
	if stoarge, ok := crm.cacheStorage.storages[id]; ok {
		return stoarge.cacheMethod, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningMethods) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}

var _ = Describe("Testthehiveapitemporatystorage", Ordered, func() {
	var (
		erm             *CacheRunningMethods
		err, errLoadEnv error
	)

	BeforeAll(func() {
		errLoadEnv = godotenv.Load("../../.env")
		fmt.Println("ERROR env:", errLoadEnv)
		fmt.Println("API KEY:", os.Getenv("GO_HIVEHOOK_THAPIKEY"))

		erm, err = New(context.Background(), 60)
	})

	Context("Тест 1. Проверяем инициализацию обработчиков", func() {
		It("Не должно быть ошибок при чтении переменных окружений", func() {
			Expect(errLoadEnv).ShouldNot(HaveOccurred())
		})

		It("Не должно быть ошибок при инициализации хранилища", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	/*
			* !!!!!!!!!!!!!!!!!
		 	* новую версию пока не выкатываю в прод потому что есть проблемы
			* с этим тестом и остальные тесты не до конца исправлены под новый
			* логгер
			* !!!!!!!!!!!!!!!!!
	*/
	a
	Context("Тест 2. Добавляем обработчики", func() {
		chZabbix := make(chan commoninterfaces.Messager)

		var listLog []simplelogger.OptionsManager
		for _, v := range []string{"error", "info"} {
			listLog = append(listLog, &confighandler.LogSet{
				MsgTypeName: v,
			})
		}
		opts := simplelogger.CreateOptions(listLog...)
		simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), "placeholder_misp", opts)
		if err != nil {
			log.Fatalf("error module 'simplelogger': %v", err)
		}

		logging := logginghandler.New(simpleLogger, chZabbix)
		logging.Start(context.Background())

		conf := confighandler.AppConfigTheHive{
			Port:   9000,
			Host:   "thehive.cloud.gcm",
			ApiKey: os.Getenv("GO_HIVEHOOK_THAPIKEY"),
		}
		apiTheHive, err := thehiveapi.New(
			logging,
			thehiveapi.WithAPIKey(conf.ApiKey),
			thehiveapi.WithHost(conf.Host),
			thehiveapi.WithPort(conf.Port))
		Expect(err).ShouldNot(HaveOccurred())

		_, err = apiTheHive.Start(context.Background())
		Expect(err).ShouldNot(HaveOccurred())

		It("Не должно быть ошибок", func() {
			chanCacheMethodResponse := make(chan CacheMethodResponse)

			erm.SetMethod("123", func() bool {
				fmt.Println("START func")

				byteBody, statusCode, err := apiTheHive.GetCaseEvent(context.Background(), "~88678416456")
				chanCacheMethodResponse <- CacheMethodResponse{
					StatusCode:   statusCode,
					Error:        err,
					BodyResponse: byteBody,
				}

				if err != nil {
					return false
				}

				if statusCode != http.StatusOK {
					return false
				}

				return true
			})

			res := <-chanCacheMethodResponse
			fmt.Printf("StatusCode: %d\n,Error:%s\n", res.StatusCode, res.Error)
		})
	})
})
