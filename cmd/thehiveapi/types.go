package thehiveapi

import (
	"encoding/json"

	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/storage"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

type TheHiveApi_New struct {
	settings theHiveApiSettings
}

// apiTheHiveModule модуль для взаимодействия с API TheHive
type apiTheHiveModule struct {
	cache            *cachingstoragewithqueue.CacheStorageWithQueue[any]
	storageCache     *storage.StorageFoundObjects
	logger           interfaces.Logger
	receivingChannel chan interfaces.ChannelRequester
	settings         theHiveApiSettings
}

// theHiveApiSettings настройки модуля
type theHiveApiSettings struct {
	apiKey             string
	host               string
	nameRegionalObject string
	cachettl           int
	port               int
}

// theHiveAPIOptions функциональные опции
type theHiveApiOptions func(*apiTheHiveModule) error

// theHiveAPIOptions_New функциональные опции
type theHiveApiOptions_New func(*TheHiveApi_New) error

type InputReguest struct {
	Data    []byte
	Command string
	RootId  string
}

// RequestChannelTheHive тип применяется для передачи запроса в модуль thehiveapi от сторонних модулей
type RequestChannelTheHive struct {
	Data       interface{}                      //какие то данные
	RequestId  string                           //UUID идентификатор запроса
	RootId     string                           //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId     string                           //идентификатор кейса в TheHive
	Command    string                           //команда
	ChanOutput chan interfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от API TheHive
type ResponseChannelTheHive struct {
	Error      error  //объект ошибки
	Data       []byte //набор данных
	Source     string //источник данных
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	StatusCode int    //статус кода ответа
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	ByteData       json.RawMessage `json:"byte_data"`           //набор данных в бинарном виде которые обрабатываются отдельно
	Service        string          `json:"service"`             //наименование сервиса
	Command        string          `json:"command"`             //команда
	RootId         string          `json:"root_id"`             //основной id, как правило это rootId case или alert
	CaseId         string          `json:"case_id"`             //id кейса
	Value          string          `json:"value"`               //устанавливаемое значение
	Username       string          `json:"username"`            //имя пользователя, необходим если нужно указать пользователя выполнившего действие
	FieldName      string          `json:"field_name"`          //некое ключевое поле
	RegionalObject string          `json:"for_regional_object"` //для кого предназначена команда
}

// Querys перечень запросов к TheHive
type Querys struct {
	Query []Query `json:"query"`
}

// Query структура запроса к TheHive
type Query struct {
	ExtraData []string `json:"extraData,omitempty"`
	Name      string   `json:"_name,omitempty"`
	IDOrName  string   `json:"idOrName,omitempty"`
	From      int64    `json:"from"`
	To        int      `json:"to,omitempty"`
}

// ErrorAnswer структура описания ошибок получаемых от TheHive
type ErrorAnswer struct {
	Err     string `json:"type"`
	Message string `json:"message"`
}
