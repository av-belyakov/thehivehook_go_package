package thehiveapi

type ErrorInformation struct {
	Err string
}

func (e ErrorInformation) Error() string {
	return e.Err
}
