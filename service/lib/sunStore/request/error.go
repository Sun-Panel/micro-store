package request

import "errors"

var ()

func DeadlyError(htppStatusCode int) error {

	switch htppStatusCode {
	case 401:
		return errors.New("unauthorized")
	case 500:
		return errors.New("server error")
	}

	return nil
}
