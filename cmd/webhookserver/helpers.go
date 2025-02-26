package webhookserver

import (
	"fmt"
	"os"

	"github.com/av-belyakov/thehivehook_go_package/internal/appname"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
)

func getInformationMessage(name, host string, port int) string {
	appStatus := fmt.Sprintf("%vproduction%v", Ansi_Bright_Blue, Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", Ansi_Bright_Red, envValue, Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), appversion.GetVersion())
	fmt.Printf("\n%v%v%s.%v\n", Bold_Font, Ansi_Bright_Green, msg, Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", Underlining, Ansi_Bright_Green, appStatus, Ansi_Reset)
	fmt.Printf("%vWebhook server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%vName regional object: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Orange, name, Ansi_Reset)
	fmt.Printf("%v  IP: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, host, Ansi_Reset)
	fmt.Printf("%v  Port: %v%d%v\n\n", Ansi_Bright_Green, Ansi_Bright_Magenta, port, Ansi_Reset)

	return msg
}
