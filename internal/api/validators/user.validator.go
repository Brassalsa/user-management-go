package validators

import "errors"

func CheckValidPassword(password string) error {

	if len(password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}
