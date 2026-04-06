package supportingfunctions

import (
	"errors"
)

// GetObjectTypeFromEventTheHive получить из TheHive события значение типа "objectType"
func GetObjectTypeFromEventTheHive(elem map[string]any) (string, error) {
	value, ok := elem["objectType"]
	if !ok {
		return "", CustomError(errors.New("the accepted object does not have the property 'objectType'"))
	}

	objType, ok := value.(string)
	if !ok {
		return "", CustomError(errors.New("it is not possible to convert a value"))
	}

	return objType, nil
}

// GetRootIdFromEventTheHive получить из TheHive события значение типа "rootId"
func GetRootIdFromEventTheHive(elem map[string]any) (string, error) {
	value, ok := elem["rootId"]
	if !ok {
		return "", CustomError(errors.New("the accepted object does not have the property 'rootId'"))
	}

	rootId, ok := value.(string)
	if !ok {
		return "", CustomError(errors.New("it is not possible to convert a value"))
	}

	return rootId, nil
}

// GetOperationFromEventTheHive получить из TheHive события значение типа "operation"
func GetOperationFromEventTheHive(elem map[string]any) (string, error) {
	value, ok := elem["operation"]
	if !ok {
		return "", CustomError(errors.New("the accepted object does not have the property 'operation'"))
	}

	operation, ok := value.(string)
	if !ok {
		return "", CustomError(errors.New("it is not possible to convert a value"))
	}

	return operation, nil
}

// GetCaseIdFromEventTheHive получить из TheHive события значение типа "caseId"
func GetCaseIdFromEventTheHive(elem map[string]any) (int, error) {
	value, ok := elem["object"]
	if !ok {
		return 0, CustomError(errors.New("the accepted object does not have the property 'object'"))
	}

	object, ok := value.(map[string]any)
	if !ok {
		return 0, CustomError(errors.New("it is not possible to convert a value"))
	}

	value, ok = object["caseId"]
	if !ok {
		return 0, CustomError(errors.New("the accepted object does not have the property 'caseId'"))
	}

	caseId, ok := value.(float64)
	if !ok {
		return 0, CustomError(errors.New("it is not possible to convert a value"))
	}

	return int(caseId), nil
}
