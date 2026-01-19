package confighandler

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func NewConfig(rootDir string) (*ConfigApp, error) {
	conf := ConfigApp{}
	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_HIVEHOOK_MAIN": "",

			//Подключение к NATS
			"GO_HIVEHOOK_NHOST":               "",
			"GO_HIVEHOOK_NPORT":               "",
			"GO_HIVEHOOK_NSUBSENDERCASE":      "",
			"GO_HIVEHOOK_NSUBSENDERALERT":     "",
			"GO_HIVEHOOK_NSUBLISTENERCOMMAND": "",

			//Подключение к TheHive
			"GO_HIVEHOOK_THHOST":   "",
			"GO_HIVEHOOK_THPORT":   "",
			"GO_HIVEHOOK_THUNAME":  "",
			"GO_HIVEHOOK_THAPIKEY": "",

			//Настройки WebHookServer
			"GO_HIVEHOOK_WEBHNAME":       "",
			"GO_HIVEHOOK_WEBHHOST":       "",
			"GO_HIVEHOOK_WEBHPORT":       "",
			"GO_HIVEHOOK_WEBHTTLTMPINFO": "",

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

		//Настройки для модуля подключения к NATS
		if viper.IsSet("NATS.host") {
			conf.AppConfigNATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			conf.AppConfigNATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cache_ttl") {
			conf.AppConfigNATS.CacheTTL = viper.GetInt("NATS.cache_ttl")
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

		//if viper.IsSet("NATS.subscribers") {
		//	nats := NATS{}
		//	if err := viper.GetViper().Unmarshal(&nats); err != nil {
		//		return err
		//	}
		//
		//	conf.AppConfigNATS.Subscribers = nats.NATS.Subscribers
		//}

		//Настройки для модуля подключения к TheHive
		if viper.IsSet("THEHIVE.host") {
			conf.AppConfigTheHive.Host = viper.GetString("THEHIVE.host")
		}
		if viper.IsSet("THEHIVE.port") {
			conf.AppConfigTheHive.Port = viper.GetInt("THEHIVE.port")
		}
		if viper.IsSet("THEHIVE.cache_ttl") {
			conf.AppConfigTheHive.CacheTTL = viper.GetInt("THEHIVE.cache_ttl")
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
		if viper.IsSet("WEBHOOKSERVER.ttl_tmp_info") {
			conf.AppConfigWebHookServer.TTLTmpInfo = viper.GetInt("WEBHOOKSERVER.ttl_tmp_info")
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

	//Настройки для модуля подключения к NATS
	if envList["GO_HIVEHOOK_NHOST"] != "" {
		conf.AppConfigNATS.Host = envList["GO_HIVEHOOK_NHOST"]
	}
	if envList["GO_HIVEHOOK_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_NPORT"]); err == nil {
			conf.AppConfigNATS.Port = p
		}
	}
	if envList["GO_HIVEHOOK_NCACHETTL"] != "" {
		if v, err := strconv.Atoi(envList["GO_HIVEHOOK_NCACHETTL"]); err == nil {
			conf.AppConfigNATS.CacheTTL = v
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

	//if envList["GO_HIVEHOOK_NSUBSCRIBERS"] != "" {
	//	subscribers := []SubscriberNATS{}
	//	if strings.Contains(envList["GO_HIVEHOOK_NSUBSCRIBERS"], ";") {
	//		tmp := strings.Split(envList["GO_HIVEHOOK_NSUBSCRIBERS"], ";")
	//		if len(tmp) > 0 {
	//			for _, v := range tmp {
	//				if subscriber, err := hundlerSubscribersString(v); err == nil {
	//					subscribers = append(subscribers, subscriber)
	//				}
	//			}
	//		}
	//	} else {
	//		if subscriber, err := hundlerSubscribersString(envList["GO_HIVEHOOK_NSUBSCRIBERS"]); err == nil {
	//			subscribers = append(subscribers, subscriber)
	//		}
	//	}
	//
	//	conf.AppConfigNATS.Subscribers = subscribers
	//}

	//Настройки для модуля подключения к TheHive
	if envList["GO_HIVEHOOK_THHOST"] != "" {
		conf.AppConfigTheHive.Host = envList["GO_HIVEHOOK_THHOST"]
	}
	if envList["GO_HIVEHOOK_THPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_THPORT"]); err == nil {
			conf.AppConfigTheHive.Port = p
		}
	}
	if envList["GO_HIVEHOOK_THCACHETTL"] != "" {
		if v, err := strconv.Atoi(envList["GO_HIVEHOOK_THCACHETTL"]); err == nil {
			conf.AppConfigTheHive.CacheTTL = v
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
	if envList["GO_HIVEHOOK_WEBHTTLTMPINFO"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_WEBHTTLTMPINFO"]); err == nil {
			conf.AppConfigWebHookServer.TTLTmpInfo = p
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

func hundlerSubscribersString(str string) (SubscriberNATS, error) {
	errMsg := "an incorrect string containing the 'subscribers' of the NATS settings was received"
	subscriber := SubscriberNATS{}

	if !strings.Contains(str, ":") {
		return subscriber, errors.New(errMsg)
	}

	tmp := strings.Split(str, ":")
	if len(tmp) == 0 {
		return subscriber, errors.New(errMsg)
	}

	responders := []string{}
	if strings.Contains(tmp[1], ",") {
		responders = append(responders, strings.Split(tmp[1], ",")...)
	} else {
		responders = append(responders, tmp[1])
	}

	subscriber.Event = tmp[0]
	subscriber.Responders = responders

	return subscriber, nil
}
