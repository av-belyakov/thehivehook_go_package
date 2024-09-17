package elasticsearchapi

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

func (h *HandlerSendData) New() error {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", h.Settings.Host, h.Settings.Port)},
		Username:  h.Settings.User,
		Password:  h.Settings.Passwd,
		Transport: &http.Transport{
			MaxIdleConns:          10,              //число открытых TCP-соединений, которые в данный момент не используются
			IdleConnTimeout:       1 * time.Second, //время, через которое закрываются такие неактивные соединения
			MaxIdleConnsPerHost:   10,              //число неактивных TCP-соединений, которые допускается устанавливать на один хост
			ResponseHeaderTimeout: 2 * time.Second, //время в течении которого сервер ожидает получение ответа после записи заголовка запроса
			DialContext: (&net.Dialer{
				Timeout: 3 * time.Second,
				//KeepAlive: 1 * time.Second,
			}).DialContext,
		},
		//RetryOnError: ,
		//RetryOnStatus: ,
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%v' %s:%d", err, f, l-1)
	}

	h.Client = es

	return nil
}

func HandlerElasticSearch(
	conf confighandler.AppConfigElasticSearch,
	logging chan<- logginghandler.MessageLogging) (*ElasticSearchModule, error) {

	module := &ElasticSearchModule{
		ChanInputModule:  make(chan SettingsInputChan),
		ChanOutputModule: make(chan SettingsOutputChan),
	}

	hsd := HandlerSendData{
		Settings: SettingsHandler{
			Port:   conf.Port,
			Host:   conf.Host,
			User:   conf.UserName,
			Passwd: conf.Passwd,
		},
	}

	if err := hsd.New(); err != nil {
		return module, err
	}

	go func() {
		for data := range module.ChanInputModule {
			switch data.Section {
			case "handling":
				index := fmt.Sprintf("%s%s", conf.Prefix, conf.Index)

				if data.Command == "add new" {
					fmt.Printf("func 'HandlerElasticSearch' Section: %s, Command: %s, Index: %s\n", data.Section, data.Command, index)
				}
			}
		}
	}()

	return module, nil
}
