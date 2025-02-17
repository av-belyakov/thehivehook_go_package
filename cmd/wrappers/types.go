package wrappers

// WrappersZabbixInteractionSettings настройки для обертки взаимодействия с модулем zabbixapi
type WrappersZabbixInteractionSettings struct {
	EventTypes  []EventType //типы событий
	NetworkHost string      //ip адрес или доменное имя
	ZabbixHost  string      //zabbix host
	NetworkPort int         //сетевой порт
}

type EventType struct {
	EventType  string
	ZabbixKey  string
	IsTransmit bool
	Handshake  Handshake
}

type Handshake struct {
	TimeInterval int
	Message      string
}
