package pkg

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (err error, hash string) {
	hashed_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return err, string(hashed_pwd)
}

func Compare(hash string, plain_text_pwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain_text_pwd))
	return err
}
