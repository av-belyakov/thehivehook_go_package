package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/internal/appname"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

func getInformationMessage(conf *confighandler.ConfigApp /*wh.name, host string, port int*/) string {
	var queryVisual string
	version, _ := appversion.GetAppVersion()
	host := "*"
	port := conf.AppConfigWebHookServer.Port
	if conf.AppConfigWebHookServer.Host != "" {
		host = conf.AppConfigWebHookServer.Host
	}

	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && (envValue == "development" || envValue == "test") {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	if ok && envValue == "test" {
		queryVisual = fmt.Sprintf("%vQuery visualization url http://%s:%d/__viz/.%v\n", constants.Ansi_Bright_Red, host, port, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf(
		"Application '%s' v%s was successfully launched",
		appname.GetName(),
		strings.Replace(version, "\n", "", -1),
	)

	fmt.Printf("\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf(
		"%v%vApplication status is '%s'%v\n",
		constants.Underlining,
		constants.Ansi_Bright_Green,
		appStatus,
		constants.Ansi_Reset,
	)
	fmt.Print(queryVisual)
	fmt.Printf(
		"%vName regional object: %v%s%v\n",
		constants.Ansi_Bright_Green,
		constants.Ansi_Bright_Orange,
		conf.AppConfigWebHookServer.Name,
		constants.Ansi_Reset,
	)
	fmt.Printf("%vWebhook server settings:%v\n", constants.Ansi_Bright_Green, constants.Ansi_Reset)
	fmt.Printf("  %vIP:%v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, host, constants.Ansi_Reset)
	fmt.Printf("  %vPort:%v%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Magenta, port, constants.Ansi_Reset)
	fmt.Printf(
		"%vConnect to NATS with address %v%s:%d%v%v, subscriptions: %v%s, %s%v\n",
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		conf.AppConfigNATS.Host,
		conf.AppConfigNATS.Port,
		constants.Ansi_Reset,
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		conf.AppConfigNATS.Subscriptions.SenderAlert,
		conf.AppConfigNATS.Subscriptions.SenderCase,
		constants.Ansi_Reset,
	)

	return msg
}
