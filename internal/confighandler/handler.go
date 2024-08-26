package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func NewConfig(rootDir string) (*ConfigApp, error) {
	conf := ConfigApp{}
	var envList map[string]string = map[string]string{
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
	}

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

		return nil
	}

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

	return &conf, nil
}
