package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/thehivehook_go_package/v2/internal/supportingfunctions"
)

func NewConfig(rootDir string) (*ConfigApp, error) {
	conf := ConfigApp{}
	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_HIVEHOOK_MAIN": "",

			//Настройки временного хранилища
			"GO_HIVEHOOK_TSTORAGEOBJECTTTL":       "",
			"GO_HIVEHOOK_TSTORAGDELAYTOSENDALERT": "",
			"GO_HIVEHOOK_TSTORAGDELAYTOSENDCASE":  "",

			//Подключение к NATS
			"GO_HIVEHOOK_NHOST":               "",
			"GO_HIVEHOOK_NPORT":               "",
			"GO_HIVEHOOK_NSUBSENDERCASE":      "",
			"GO_HIVEHOOK_NSUBSENDERALERT":     "",
			"GO_HIVEHOOK_NSUBLISTENERCOMMAND": "",

			//Подключение к TheHive
			"GO_HIVEHOOK_THHOST":   "",
			"GO_HIVEHOOK_THPORT":   "",
			"GO_HIVEHOOK_THAPIKEY": "",

			//Настройки WebHookServer
			"GO_HIVEHOOK_WEBHNAME": "",
			"GO_HIVEHOOK_WEBHHOST": "",
			"GO_HIVEHOOK_WEBHPORT": "",

			//Настройки доступа к БД в которую будут записыватся логи
			"GO_HIVEHOOK_DBWLOGHOST":        "",
			"GO_HIVEHOOK_DBWLOGPORT":        "",
			"GO_HIVEHOOK_DBWLOGNAME":        "",
			"GO_HIVEHOOK_DBWLOGUSER":        "",
			"GO_HIVEHOOK_DBWLOGPASSWD":      "",
			"GO_HIVEHOOK_DBWLOGSTORAGENAME": "",
		}
	)

	getFileName := func(sf, confPath string, lfs []fs.DirEntry) (string, error) {
		for _, v := range lfs {
			if v.Name() == sf && !v.IsDir() {
				return path.Join(confPath, v.Name()), nil
			}
		}

		return "", fmt.Errorf("file '%s' is not found", sf)
	}

	setCommonSettings := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			conf.CommonAppConfig.LogList = ls.Logging
		}

		z := ZabbixSet{}
		if ok := viper.IsSet("ZABBIX"); ok {
			if err := viper.GetViper().Unmarshal(&z); err != nil {
				return err
			}

			np := 10051
			if z.Zabbix.NetworkPort != 0 && z.Zabbix.NetworkPort < 65536 {
				np = z.Zabbix.NetworkPort
			}

			conf.CommonAppConfig.Zabbix = ZabbixOptions{
				NetworkPort: np,
				NetworkHost: z.Zabbix.NetworkHost,
				ZabbixHost:  z.Zabbix.ZabbixHost,
				EventTypes:  z.Zabbix.EventTypes,
			}
		}

		return nil
	}

	setSpecial := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		//Общие настройки конфигурационного файла
		if viper.IsSet("COMMONINFO.file_name") {
			conf.CommonInfo.FileName = viper.GetString("COMMONINFO.file_name")
		}

		//Настройки временного хранилища
		if viper.IsSet("TEMPORARYSTORAGE.storage_object_ttl") {
			conf.AppConfigTemporaryStorage.StorageObjectTTL = viper.GetInt("TEMPORARYSTORAGE.storage_object_ttl")
		}
		if viper.IsSet("TEMPORARYSTORAGE.storage_delay_to_send_alert") {
			conf.AppConfigTemporaryStorage.StorageDelayToSendAlert = viper.GetInt("TEMPORARYSTORAGE.storage_delay_to_send_alert")
		}
		if viper.IsSet("TEMPORARYSTORAGE.storage_delay_to_send_case") {
			conf.AppConfigTemporaryStorage.StorageDelayToSendCase = viper.GetInt("TEMPORARYSTORAGE.storage_delay_to_send_case")
		}

		//Настройки для модуля подключения к NATS
		if viper.IsSet("NATS.host") {
			conf.AppConfigNATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			conf.AppConfigNATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.subscriptions.sender_case") {
			conf.AppConfigNATS.Subscriptions.SenderCase = viper.GetString("NATS.subscriptions.sender_case")
		}
		if viper.IsSet("NATS.subscriptions.sender_alert") {
			conf.AppConfigNATS.Subscriptions.SenderAlert = viper.GetString("NATS.subscriptions.sender_alert")
		}
		if viper.IsSet("NATS.subscriptions.listener_command") {
			conf.AppConfigNATS.Subscriptions.ListenerCommand = viper.GetString("NATS.subscriptions.listener_command")
		}

		//Настройки для модуля подключения к TheHive
		if viper.IsSet("THEHIVE.host") {
			conf.AppConfigTheHive.Host = viper.GetString("THEHIVE.host")
		}
		if viper.IsSet("THEHIVE.port") {
			conf.AppConfigTheHive.Port = viper.GetInt("THEHIVE.port")
		}
		if viper.IsSet("THEHIVE.api_key") {
			conf.AppConfigTheHive.ApiKey = viper.GetString("THEHIVE.api_key")
		}

		//	Настройки основного API сервера
		if viper.IsSet("WEBHOOKSERVER.name") {
			conf.AppConfigWebHookServer.Name = viper.GetString("WEBHOOKSERVER.name")
		}
		if viper.IsSet("WEBHOOKSERVER.host") {
			conf.AppConfigWebHookServer.Host = viper.GetString("WEBHOOKSERVER.host")
		}
		if viper.IsSet("WEBHOOKSERVER.port") {
			conf.AppConfigWebHookServer.Port = viper.GetInt("WEBHOOKSERVER.port")
		}

		//Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("DATABASEWRITELOG.host") {
			conf.AppConfigWriteLogDB.Host = viper.GetString("DATABASEWRITELOG.host")
		}
		if viper.IsSet("DATABASEWRITELOG.port") {
			conf.AppConfigWriteLogDB.Port = viper.GetInt("DATABASEWRITELOG.port")
		}
		if viper.IsSet("DATABASEWRITELOG.user") {
			conf.AppConfigWriteLogDB.User = viper.GetString("DATABASEWRITELOG.user")
		}
		if viper.IsSet("DATABASEWRITELOG.namedb") {
			conf.AppConfigWriteLogDB.NameDB = viper.GetString("DATABASEWRITELOG.namedb")
		}
		if viper.IsSet("DATABASEWRITELOG.storage_name_db") {
			conf.AppConfigWriteLogDB.StorageNameDB = viper.GetString("DATABASEWRITELOG.storage_name_db")
		}

		return nil
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return &conf, err
	}

	confPath := path.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return &conf, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return &conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return &conf, err
	}

	var fn string
	switch envList["GO_HIVEHOOK_MAIN"] {
	case "development":
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return &conf, err
		}
	case "test":
		fn, err = getFileName("config_test.yml", confPath, list)
		if err != nil {
			return &conf, err
		}
	default:
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return &conf, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return &conf, err
	}

	//Настройки временного хранилища
	if envList["GO_HIVEHOOK_TSTORAGEOBJECTTTL"] != "" {
		if v, err := strconv.Atoi(envList["GO_HIVEHOOK_TSTORAGEOBJECTTTL"]); err == nil {
			conf.AppConfigTemporaryStorage.StorageObjectTTL = v
		}
	}
	if envList["GO_HIVEHOOK_TSTORAGDELAYTOSENDALERT"] != "" {
		if v, err := strconv.Atoi(envList["GO_HIVEHOOK_TSTORAGDELAYTOSENDALERT"]); err == nil {
			conf.AppConfigTemporaryStorage.StorageDelayToSendAlert = v
		}
	}
	if envList["GO_HIVEHOOK_TSTORAGDELAYTOSENDCASE"] != "" {
		if v, err := strconv.Atoi(envList["GO_HIVEHOOK_TSTORAGDELAYTOSENDCASE"]); err == nil {
			conf.AppConfigTemporaryStorage.StorageDelayToSendCase = v
		}
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_HIVEHOOK_NHOST"] != "" {
		conf.AppConfigNATS.Host = envList["GO_HIVEHOOK_NHOST"]
	}
	if envList["GO_HIVEHOOK_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_NPORT"]); err == nil {
			conf.AppConfigNATS.Port = p
		}
	}
	if envList["GO_HIVEHOOK_NSUBSENDERCASE"] != "" {
		conf.AppConfigNATS.Subscriptions.SenderCase = envList["GO_HIVEHOOK_NSUBSENDERCASE"]
	}
	if envList["GO_HIVEHOOK_NSUBSENDERALERT"] != "" {
		conf.AppConfigNATS.Subscriptions.SenderAlert = envList["GO_HIVEHOOK_NSUBSENDERALERT"]
	}
	if envList["GO_HIVEHOOK_NSUBLISTENERCOMMAND"] != "" {
		conf.AppConfigNATS.Subscriptions.ListenerCommand = envList["GO_HIVEHOOK_NSUBLISTENERCOMMAND"]
	}

	//Настройки для модуля подключения к TheHive
	if envList["GO_HIVEHOOK_THHOST"] != "" {
		conf.AppConfigTheHive.Host = envList["GO_HIVEHOOK_THHOST"]
	}
	if envList["GO_HIVEHOOK_THPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_THPORT"]); err == nil {
			conf.AppConfigTheHive.Port = p
		}
	}
	if envList["GO_HIVEHOOK_THAPIKEY"] != "" {
		conf.AppConfigTheHive.ApiKey = envList["GO_HIVEHOOK_THAPIKEY"]
	}

	//Настройки основного API сервера
	if envList["GO_HIVEHOOK_WEBHNAME"] != "" {
		conf.AppConfigWebHookServer.Name = envList["GO_HIVEHOOK_WEBHNAME"]
	}
	if envList["GO_HIVEHOOK_WEBHHOST"] != "" {
		conf.AppConfigWebHookServer.Host = envList["GO_HIVEHOOK_WEBHHOST"]
	}
	if envList["GO_HIVEHOOK_WEBHPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_WEBHPORT"]); err == nil {
			conf.AppConfigWebHookServer.Port = p
		}
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_HIVEHOOK_DBWLOGHOST"] != "" {
		conf.AppConfigWriteLogDB.Host = envList["GO_HIVEHOOK_DBWLOGHOST"]
	}
	if envList["GO_HIVEHOOK_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_DBWLOGPORT"]); err == nil {
			conf.AppConfigWriteLogDB.Port = p
		}
	}
	if envList["GO_HIVEHOOK_DBWLOGNAME"] != "" {
		conf.AppConfigWriteLogDB.NameDB = envList["GO_HIVEHOOK_DBWLOGNAME"]
	}
	if envList["GO_HIVEHOOK_DBWLOGUSER"] != "" {
		conf.AppConfigWriteLogDB.User = envList["GO_HIVEHOOK_DBWLOGUSER"]
	}
	if envList["GO_HIVEHOOK_DBWLOGPASSWD"] != "" {
		conf.AppConfigWriteLogDB.Passwd = envList["GO_HIVEHOOK_DBWLOGPASSWD"]
	}
	if envList["GO_HIVEHOOK_DBWLOGSTORAGENAME"] != "" {
		conf.AppConfigWriteLogDB.StorageNameDB = envList["GO_HIVEHOOK_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return &conf, err
	}

	return &conf, nil
}
