package webhookserver

import (
	"fmt"
	"os"

	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/internal/appname"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
)

func getInformationMessage(name, host string, port int) string {
	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), appversion.GetVersion())
	fmt.Printf("\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", constants.Underlining, constants.Ansi_Bright_Green, appStatus, constants.Ansi_Reset)
	fmt.Printf("%vWebhook server settings:%v\n", constants.Ansi_Bright_Green, constants.Ansi_Reset)
	fmt.Printf("%vName regional object: %v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Orange, name, constants.Ansi_Reset)
	fmt.Printf("%v  IP: %v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, host, constants.Ansi_Reset)
	fmt.Printf("%v  Port: %v%d%v\n\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Magenta, port, constants.Ansi_Reset)

	return msg
}
