package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

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
			"GO_HIVEHOOK_NHOST":        "",
			"GO_HIVEHOOK_NPORT":        "",
			"GO_HIVEHOOK_SUBJECTCASE":  "",
			"GO_HIVEHOOK_SUBJECTALERT": "",

			//Подключение к TheHive
			"GO_HIVEHOOK_THHOST":   "",
			"GO_HIVEHOOK_THPORT":   "",
			"GO_HIVEHOOK_THUNAME":  "",
			"GO_HIVEHOOK_THAPIKEY": "",

			// Подключение к СУБД Elasticsearch
			"GO_HIVEHOOK_ESHOST":   "",
			"GO_HIVEHOOK_ESPORT":   "",
			"GO_HIVEHOOK_ESUSER":   "",
			"GO_HIVEHOOK_ESPASSWD": "",
			"GO_HIVEHOOK_ESPREFIX": "",
			"GO_HIVEHOOK_ESINDEX":  "",

			//Настройки основного API серврера
			"GO_HIVEHOOK_HHOST": "",
			"GO_HIVEHOOK_HPORT": "",
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
		viper.SetConfigType("yaml")
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
		viper.SetConfigType("yaml")
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
		if viper.IsSet("NATS.subject_case") {
			conf.AppConfigNATS.SubjectCase = viper.GetString("NATS.subject_case")
		}
		if viper.IsSet("NATS.subject_alert") {
			conf.AppConfigNATS.SubjectAlert = viper.GetString("NATS.subject_alert")
		}

		//Настройки для модуля подключения к TheHive
		if viper.IsSet("THEHIVE.host") {
			conf.AppConfigTheHive.Host = viper.GetString("THEHIVE.host")
		}
		if viper.IsSet("THEHIVE.port") {
			conf.AppConfigTheHive.Port = viper.GetInt("THEHIVE.port")
		}
		if viper.IsSet("THEHIVE.user_name") {
			conf.AppConfigTheHive.UserName = viper.GetString("THEHIVE.user_name")
		}
		if viper.IsSet("THEHIVE.api_key") {
			conf.AppConfigTheHive.ApiKey = viper.GetString("THEHIVE.api_key")
		}

		// Настройки для модуля подключения к СУБД ElasticSearch
		if viper.IsSet("ELASTICSEARCH.host") {
			conf.AppConfigElasticSearch.Host = viper.GetString("ELASTICSEARCH.host")
		}
		if viper.IsSet("ELASTICSEARCH.port") {
			conf.AppConfigElasticSearch.Port = viper.GetInt("ELASTICSEARCH.port")
		}
		if viper.IsSet("ELASTICSEARCH.user") {
			conf.AppConfigElasticSearch.UserName = viper.GetString("ELASTICSEARCH.user")
		}
		if viper.IsSet("ELASTICSEARCH.prefix") {
			conf.AppConfigElasticSearch.Prefix = viper.GetString("ELASTICSEARCH.prefix")
		}
		if viper.IsSet("ELASTICSEARCH.index") {
			conf.AppConfigElasticSearch.Index = viper.GetString("ELASTICSEARCH.index")
		}

		//	Настройки основного API сервера
		if viper.IsSet("HOOKSERVER.host") {
			conf.AppConfigHookServer.Host = viper.GetString("HOOKSERVER.host")
		}
		if viper.IsSet("HOOKSERVER.port") {
			conf.AppConfigHookServer.Port = viper.GetInt("HOOKSERVER.port")
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

	confPath := path.Join(rootPath, "configs")

	list, err := os.ReadDir(confPath)
	if err != nil {
		return &conf, err
	}

	fileNameCommon, err := getFileName("config.yaml", confPath, list)
	if err != nil {
		return &conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return &conf, err
	}

	var fn string
	if envList["GO_HIVEHOOK_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yaml", confPath, list)
		if err != nil {
			return &conf, err
		}
	} else {
		fn, err = getFileName("config_prod.yaml", confPath, list)
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
	if envList["GO_HIVEHOOK_SUBJECTCASE"] != "" {
		conf.AppConfigNATS.SubjectCase = envList["GO_HIVEHOOK_SUBJECTCASE"]
	}
	if envList["GO_HIVEHOOK_SUBJECTALERT"] != "" {
		conf.AppConfigNATS.SubjectAlert = envList["GO_HIVEHOOK_SUBJECTALERT"]
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
	if envList["GO_HIVEHOOK_THUNAME"] != "" {
		conf.AppConfigTheHive.UserName = envList["GO_HIVEHOOK_THUNAME"]
	}
	if envList["GO_HIVEHOOK_THAPIKEY"] != "" {
		conf.AppConfigTheHive.ApiKey = envList["GO_HIVEHOOK_THAPIKEY"]
	}

	//Настройки для модуля подключения к СУБД ElasticSearch
	if envList["GO_HIVEHOOK_ESHOST"] != "" {
		conf.AppConfigElasticSearch.Host = envList["GO_HIVEHOOK_ESHOST"]
	}
	if envList["GO_HIVEHOOK_ESPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_ESPORT"]); err == nil {
			conf.AppConfigElasticSearch.Port = p
		}
	}
	if envList["GO_HIVEHOOK_ESUSER"] != "" {
		conf.AppConfigElasticSearch.UserName = envList["GO_HIVEHOOK_ESUSER"]
	}
	if envList["GO_HIVEHOOK_ESPASSWD"] != "" {
		conf.AppConfigElasticSearch.Passwd = envList["GO_HIVEHOOK_ESPASSWD"]
	}
	if envList["GO_HIVEHOOK_ESPREFIX"] != "" {
		conf.AppConfigElasticSearch.Prefix = envList["GO_HIVEHOOK_ESPREFIX"]
	}
	if envList["GO_HIVEHOOK_ESINDEX"] != "" {
		conf.AppConfigElasticSearch.Index = envList["GO_HIVEHOOK_ESINDEX"]
	}

	//Настройки основного API сервера
	if envList["GO_HIVEHOOK_HHOST"] != "" {
		conf.AppConfigHookServer.Host = envList["GO_HIVEHOOK_HHOST"]
	}
	if envList["GO_HIVEHOOK_HPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_HIVEHOOK_HPORT"]); err == nil {
			conf.AppConfigHookServer.Port = p
		}
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return &conf, err
	}

	return &conf, nil
}
