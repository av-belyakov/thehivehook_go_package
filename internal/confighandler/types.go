package confighandler

type ConfigApp struct {
	CommonInfo
	CommonAppConfig
	AppConfigNATS
	AppConfigTheHive
	AppConfigElasticSearch
	AppConfigHookServer
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
	Host        string           `validate:"required" yaml:"host"`
	Subscribers []SubscriberNATS `yaml:"subscribers"`
}

type AppConfigTheHive struct {
	Port     int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host     string `validate:"required" yaml:"host"`
	UserName string `validate:"required" yaml:"user_name"`
	ApiKey   string `validate:"required"`
}

type AppConfigElasticSearch struct {
	Port     int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host     string `validate:"required" yaml:"host"`
	UserName string `validate:"required" yaml:"user"`
	Passwd   string `validate:"required"`
	Prefix   string `yaml:"prefix"`
	Index    string `validate:"required" yaml:"index"`
}

type AppConfigHookServer struct {
	Port int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host string `validate:"required" yaml:"host"`
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
