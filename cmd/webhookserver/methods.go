package webhookserver

import "fmt"

func (e *CreateCaseError) Error() string {
	return fmt.Sprintf("%s:%v", e.Type, e.Err)
}

func (e *CreateCaseError) Is(target error) bool {
	ccerr, ok := target.(*CreateCaseError)
	if !ok {
		return false
	}

	return e.Type == ccerr.Type
}
