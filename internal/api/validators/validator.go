package validators

import "errors"

func CheckValidPassword(password string) error {

	if len(password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}

func CheckValidRole(role string) error {
	if role != "user" && role != "admin" {
		return errors.New("role is not valid")
	}
	return nil
}
