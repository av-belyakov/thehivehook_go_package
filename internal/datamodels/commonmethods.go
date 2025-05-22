package datamodels

import "fmt"

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s:%v", e.Type, e.Err)
}

func (e *CustomError) Is(target error) bool {
	ccerr, ok := target.(*CustomError)
	if !ok {
		return false
	}

	return e.Type == ccerr.Type
}
