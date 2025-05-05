package webhookserver

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/internal/appname"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func getInformationMessage(name, host string, port int) string {
	version, _ := appversion.GetAppVersion()

	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && (envValue == "development" || envValue == "test") {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), strings.Replace(version, "\n", "", -1))
	fmt.Printf("\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", constants.Underlining, constants.Ansi_Bright_Green, appStatus, constants.Ansi_Reset)
	fmt.Printf("%vWebhook server settings:%v\n", constants.Ansi_Bright_Green, constants.Ansi_Reset)
	fmt.Printf("%vName regional object: %v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Orange, name, constants.Ansi_Reset)
	fmt.Printf("%v  IP: %v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, host, constants.Ansi_Reset)
	fmt.Printf("%v  Port: %v%d%v\n\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Magenta, port, constants.Ansi_Reset)

	return msg
}

// GetObjectType получить из карты значение типа "objectType"
func GetObjectType(elem map[string]any) (string, error) {
	value, ok := elem["objectType"]
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("the accepted object does not have the property 'objectType'"))
	}

	objType, ok := value.(string)
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("it is not possible to convert a value"))
	}

	return objType, nil
}

// GetRootId получить из карты значение типа "rootId"
func GetRootId(elem map[string]any) (string, error) {
	value, ok := elem["rootId"]
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("the accepted object does not have the property 'rootId'"))
	}

	rootId, ok := value.(string)
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("it is not possible to convert a value"))
	}

	return rootId, nil
}

// GetOperation получить из карты значение типа "operation"
func GetOperation(elem map[string]any) (string, error) {
	value, ok := elem["operation"]
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("the accepted object does not have the property 'operation'"))
	}

	operation, ok := value.(string)
	if !ok {
		return "", supportingfunctions.CustomError(errors.New("it is not possible to convert a value"))
	}

	return operation, nil
}

// GetCaseId получить из карты значение типа "caseId"
func GetCaseId(elem map[string]any) (int, error) {
	value, ok := elem["object"]
	if !ok {
		return 0, supportingfunctions.CustomError(errors.New("the accepted object does not have the property 'object'"))
	}

	object, ok := value.(map[string]any)
	if !ok {
		return 0, supportingfunctions.CustomError(errors.New("it is not possible to convert a value"))
	}

	value, ok = object["caseId"]
	if !ok {
		return 0, supportingfunctions.CustomError(errors.New("the accepted object does not have the property 'caseId'"))
	}

	caseId, ok := value.(float64)
	if !ok {
		return 0, supportingfunctions.CustomError(errors.New("it is not possible to convert a value"))
	}

	return int(caseId), nil
}
