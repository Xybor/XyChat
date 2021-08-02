package tokens

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xybor/xychat/helpers"
)

type UserToken struct {
	ID         uint
	Expiration time.Duration
}

type UserTokenClaims struct {
	ID uint `json:"id"`
	*jwt.StandardClaims
}

func getSecret(token *jwt.Token) (interface{}, error) {
	secret_string := helpers.ReadEnv("token_secret", "1ns3cur3t0k3n")

	secret := []byte(secret_string)
	return secret, nil
}

func (ut UserToken) Create() (string, error) {
	claims := UserTokenClaims{
		ID: ut.ID,
		StandardClaims: &jwt.StandardClaims{
			Audience:  "Xybor",
			ExpiresAt: time.Now().Add(ut.Expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Xychat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := getSecret(&jwt.Token{})
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(secret)
	return signedToken, err
}

func (ut *UserToken) Validate(signedToken string) error {
	claims := UserTokenClaims{}
	_, err := jwt.ParseWithClaims(
		signedToken,
		&claims,
		getSecret,
	)

	if err != nil {
		log.Print(err)
		return errors.New("invalid token")
	}

	expiration := claims.StandardClaims.ExpiresAt
	if expiration-time.Now().Unix() < 0 {
		return errors.New("token is expired")
	}

	ut.ID = claims.ID

	return nil
}
