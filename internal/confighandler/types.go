package confighandler

type ConfigApp struct {
	CommonInfo
	CommonAppConfig
	AppConfigNATS
	AppConfigTheHive
	AppConfigWebHookServer
}

type CommonInfo struct {
	FileName string `validate:"required" yaml:"filename"`
}

type CommonAppConfig struct {
	LogList []LogSet
	Zabbix  ZabbixOptions
}

type Logs struct {
	Logging []LogSet
}

type LogSet struct {
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
}

type ZabbixSet struct {
	Zabbix ZabbixOptions
}

type ZabbixOptions struct {
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	EventTypes  []EventType `yaml:"eventType"`
}

type EventType struct {
	IsTransmit bool      `yaml:"isTransmit"`
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
}

type Handshake struct {
	TimeInterval int    `yaml:"timeInterval"`
	Message      string `validate:"required" yaml:"message"`
}

type AppConfigNATS struct {
	Port        int              `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL    int              `validate:"gt=10,lte=86400" yaml:"cacheTtl"`
	Host        string           `validate:"required" yaml:"host"`
	Subscribers []SubscriberNATS `yaml:"subscribers"`
}

type AppConfigTheHive struct {
	Port     int    `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL int    `validate:"gt=10,lte=86400" yaml:"cacheTtl"`
	Host     string `validate:"required" yaml:"host"`
	ApiKey   string `validate:"required"`
}

type AppConfigWebHookServer struct {
	TTLTmpInfo int    `validate:"gt=9,lte=86400" yaml:"ttlTmpInfo"`
	Port       int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host       string `validate:"required" yaml:"host"`
	Name       string `validate:"required" yaml:"name"`
}

type NATS struct {
	NATS SubscribersNATS
}

type SubscribersNATS struct {
	Subscribers []SubscriberNATS `yaml:"subscribers"`
}

type SubscriberNATS struct {
	Event      string   `validate:"required" yaml:"event"`
	Responders []string `yaml:"responders"`
}
