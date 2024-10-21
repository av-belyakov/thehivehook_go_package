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

// GetApplicationWebHookServer настройки основного сервера API
func (conf *ConfigApp) GetApplicationWebHookServer() *AppConfigWebHookServer {
	return &conf.AppConfigWebHookServer
}

// GetApplicationSqlite настройки подключения к SQLite
func (conf *ConfigApp) GetApplicationSqlite() *AppConfigSqlite {
	return &conf.AppConfigSqlite
}

// Clean
func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}
