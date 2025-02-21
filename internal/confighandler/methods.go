package confighandler

import "errors"

func (conf *ConfigApp) GetCommonInfo() *CommonInfo {
	return &conf.CommonInfo
}

// GetCommonApplication общие настройки приложения
func (conf *ConfigApp) GetCommonApplication() *CommonAppConfig {
	return &conf.CommonAppConfig
}

// GetListLogs настройки логирования
func (conf *ConfigApp) GetListLogs() []*LogSet {
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

// GetApplicationWriteLogDB настройки доступа к БД для логирования данных
func (conf *ConfigApp) GetApplicationWriteLogDB() *AppConfigWriteLogDB {
	return &conf.AppConfigWriteLogDB
}

// Clean
func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}

// SetNameMessageType наименование тпа логирования
func (l *LogSet) SetNameMessageType(v string) error {
	if v == "" {
		return errors.New("the value 'MsgTypeName' must not be empty")
	}

	return nil
}

// SetMaxLogFileSize максимальный размер файла для логирования
func (l *LogSet) SetMaxLogFileSize(v int) error {
	if v < 1000 {
		return errors.New("the value 'MaxFileSize' must not be less than 1000")
	}

	return nil
}

// SetPathDirectory путь к директории логирования
func (l *LogSet) SetPathDirectory(v string) error {
	if v == "" {
		return errors.New("the value 'PathDirectory' must not be empty")
	}

	return nil
}

// SetWritingStdout запись логов на вывод stdout
func (l *LogSet) SetWritingStdout(v bool) {
	l.WritingStdout = v
}

// SetWritingFile запись логов в файл
func (l *LogSet) SetWritingFile(v bool) {
	l.WritingFile = v
}

// SetWritingDB запись логов  в БД
func (l *LogSet) SetWritingDB(v bool) {
	l.WritingDB = v
}

// GetNameMessageType наименование тпа логирования
func (l *LogSet) GetNameMessageType() string {
	return l.MsgTypeName
}

// GetMaxLogFileSize максимальный размер файла для логирования
func (l *LogSet) GetMaxLogFileSize() int {
	return l.MaxFileSize
}

// GetPathDirectory путь к директории логирования
func (l *LogSet) GetPathDirectory() string {
	return l.PathDirectory
}

// GetWritingStdout запись логов на вывод stdout
func (l *LogSet) GetWritingStdout() bool {
	return l.WritingStdout
}

// GetWritingFile запись логов в файл
func (l *LogSet) GetWritingFile() bool {
	return l.WritingFile
}

// GetWritingDB запись логов  в БД
func (l *LogSet) GetWritingDB() bool {
	return l.WritingDB
}
