package tokens

import (
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

type Tokens struct {
	key []byte
}

func New(path string) (*Tokens, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %+v", err)
	}

	tokens := Tokens{
		key: data,
	}

	return &tokens, nil
}

func (t *Tokens) NewToken(id string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": id,
		},
	)

	tokenString, err := token.SignedString(t.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %+v", err)
	}

	return tokenString, nil
}
