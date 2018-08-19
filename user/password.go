package user

import (
	"fmt"
	"crypto/sha256"
	"encoding/base64"
)

type Password interface {

}

func hashPassword(pwd string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(pwd))
	if err != nil {
		return "", err
	}
	hash := h.Sum(nil)
	strHash := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprint(strHash), nil
}
