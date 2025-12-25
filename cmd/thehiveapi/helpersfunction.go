package thehiveapi

import (
	"encoding/json"
	"errors"
)

func getRequestCommandData(i any) (RequestCommand, error) {
	rc := RequestCommand{}
	data, ok := i.([]byte)
	if !ok {
		return rc, errors.New("'it is not possible to convert a some value to a []byte'")
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return rc, err
	}

	return rc, nil
}
