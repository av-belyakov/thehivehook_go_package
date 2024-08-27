package natsapi

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var (
	ns    *natsStorage
	once  sync.Once
	mnats ModuleNATS
)

type natsStorage struct {
	storage map[string]messageDescriptors
	sync.Mutex
}

type messageDescriptors struct {
	timeCreate int64
	msgNats    *nats.Msg
}

func NewStorageNATS() *natsStorage {
	once.Do(func() {
		ns = &natsStorage{storage: make(map[string]messageDescriptors)}

		go checkLiveTime(ns)
	})

	return ns
}

func checkLiveTime(ns *natsStorage) {
	for range time.Tick(5 * time.Second) {
		go func() {
			for k, v := range ns.storage {
				if time.Now().Unix() > (v.timeCreate + 360) {
					ns.deleteElement(k)
				}
			}
		}()
	}
}

func (ns *natsStorage) setElement(m *nats.Msg) string {
	id := uuid.New().String()

	ns.Lock()
	defer ns.Unlock()

	ns.storage[id] = messageDescriptors{
		timeCreate: time.Now().Unix(),
		msgNats:    m,
	}

	return id
}

func (ns *natsStorage) deleteElement(id string) {
	ns.Lock()
	defer ns.Unlock()

	delete(ns.storage, id)
}

// NewClientNATS создает новое подключение к NATS
func NewClientNATS(
	conf confighandler.AppConfigNATS,
	logging chan<- logginghandler.MessageLogging) (*ModuleNATS, error) {

	mnats.chanOutputNATS = make(chan SettingsOutputChan)
	//инициируем хранилище для дескрипторов сообщений NATS
	ns = NewStorageNATS()

	if conf.SubjectCase == "" && conf.SubjectAlert == "" {
		_, f, l, _ := runtime.Caller(0)
		return &mnats, fmt.Errorf("'there is not a single subscription available for NATS in the configuration file' %s:%d", f, l-1)
	}

	subjects := map[string]string{
		"subject_case":  conf.SubjectCase,
		"subject_alert": conf.SubjectAlert,
	}

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		return &mnats, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-4)
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logging <- logginghandler.MessageLogging{
			MsgData: fmt.Sprintf("the connection with NATS has been disconnected %s:%d", f, l-4),
			MsgType: "error",
		}
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		logging <- logginghandler.MessageLogging{
			MsgData: fmt.Sprintf("the connection to NATS has been re-established %s:%d", f, l-4),
			MsgType: "info",
		}
	})

	for k, v := range subjects {
		//не добавляем обработчик если подписка пуста
		if v == "" {
			continue
		}

		nc.Subscribe(v, func(m *nats.Msg) {
			mnats.chanOutputNATS <- SettingsOutputChan{
				MsgId:       ns.setElement(m),
				SubjectType: k,
				Data:        m.Data,
			}
		})

	}

	log.Printf("Connect to NATS with address %s:%d\n", conf.Host, conf.Port)

	return &mnats, nil
}
