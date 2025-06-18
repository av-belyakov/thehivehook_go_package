package webhookserver

import (
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

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
