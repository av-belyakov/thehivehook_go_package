package confighandler

// ConfigApp настройки приложения
type ConfigApp struct {
	CommonInfo
	CommonAppConfig
	AppConfigNATS
	AppConfigTheHive
	AppConfigWebHookServer
	AppConfigWriteLogDB
}

// CommonInfo общая информация
type CommonInfo struct {
	FileName string `validate:"required" yaml:"filename"`
}

// CommonAppConfig общие настройки
type CommonAppConfig struct {
	LogList []*LogSet
	Zabbix  ZabbixOptions
}

type Logs struct {
	Logging []*LogSet
}

type LogSet struct {
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
}

type ZabbixSet struct {
	Zabbix ZabbixOptions
}

type ZabbixOptions struct {
	EventTypes  []EventType `yaml:"eventType"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
}

type EventType struct {
	Handshake  Handshake `yaml:"handshake"`
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	IsTransmit bool      `yaml:"isTransmit"`
}

type Handshake struct {
	Message      string `validate:"required" yaml:"message"`
	TimeInterval int    `yaml:"timeInterval"`
}

type AppConfigNATS struct {
	Subscriptions SubscriptionsNATS `yaml:"subscriptions"`
	Host          string            `validate:"required" yaml:"host"`
	Port          int               `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL      int               `validate:"gt=10,lte=86400" yaml:"cache_ttl"`
}

type AppConfigTheHive struct {
	ApiKey   string `validate:"required"`
	Host     string `validate:"required" yaml:"host"`
	Port     int    `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL int    `validate:"gt=10,lte=86400" yaml:"cache_ttl"`
}

type AppConfigWebHookServer struct {
	Host       string `validate:"required" yaml:"host"`
	Name       string `validate:"required" yaml:"name"`
	TTLTmpInfo int    `validate:"gt=9,lte=86400" yaml:"ttl_tmp_info"`
	Port       int    `validate:"gt=0,lte=65535" yaml:"port"`
}

type SubscriptionsNATS struct {
	SenderCase      string `validate:"required" yaml:"sender_case"`
	SenderAlert     string `validate:"required" yaml:"sender_alert"`
	ListenerCommand string `validate:"required" yaml:"listener_command"`
}

type SubscriberNATS struct {
	Responders []string `yaml:"responders"`
	Event      string   `validate:"required" yaml:"event"`
}

type AppConfigWriteLogDB struct {
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Passwd        string `yaml:"passwd"`
	NameDB        string `yaml:"namedb"`
	StorageNameDB string `yaml:"storage_name_db"`
	Port          int    `validate:"gt=0,lte=65535" yaml:"port"`
}
