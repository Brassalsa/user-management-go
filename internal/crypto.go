package internal

import "golang.org/x/crypto/bcrypt"

func HashString(str *string) error {
	res, err := bcrypt.GenerateFromPassword([]byte(*str), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	*str = string(res)
	return nil
}

func CompareHash(hash, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
