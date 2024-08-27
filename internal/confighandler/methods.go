package confighandler

func (conf *ConfigApp) GetCommonInfo() *CommonInfo {
	return &conf.CommonInfo
}

// GetCommonApplication общие настройки приложения
func (conf *ConfigApp) GetCommonApplication() *CommonAppConfig {
	return &conf.CommonAppConfig
}

// GetListLogs настройки логирования
func (conf *ConfigApp) GetListLogs() []LogSet {
	return conf.LogList
}

// GetApplicationNATS настройки взаимодействия с NATS
func (conf *ConfigApp) GetApplicationNATS() *AppConfigNATS {
	return &conf.AppConfigNATS
}

// GetApplicationTheHive настройки взаимодействия с TheHive
func (conf *ConfigApp) GetApplicationTheHive() *AppConfigTheHive {
	return &conf.AppConfigTheHive
}

// GetApplicationElasticsearch настройки взаимодействия с Elasticsearch
func (conf *ConfigApp) GetApplicationElasticsearch() *AppConfigElasticSearch {
	return &conf.AppConfigElasticSearch
}

// GetApplicationHookServer настройки основного сервера API
func (conf *ConfigApp) GetApplicationHookServer() *AppConfigHookServer {
	return &conf.AppConfigHookServer
}

// Clean
func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}
