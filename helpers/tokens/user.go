package tokens

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xybor/xychat/helpers"
)

type userToken struct {
	uid        uint
	expiration time.Duration
}

type userTokenClaims struct {
	UID uint `json:"id"`
	*jwt.StandardClaims
}

func getSecret(token *jwt.Token) (interface{}, error) {
	secret_string := helpers.ReadEnv("token_secret", "1ns3cur3t0k3n")

	secret := []byte(secret_string)
	return secret, nil
}

func CreateUserToken(id uint, expiration time.Duration) userToken {
	return userToken{uid: id, expiration: expiration}
}

func (ut userToken) GetUID() uint {
	return ut.uid
}

func (ut userToken) Generate() (string, error) {
	claims := userTokenClaims{
		UID: ut.uid,
		StandardClaims: &jwt.StandardClaims{
			Audience:  "Xybor",
			ExpiresAt: time.Now().Add(ut.expiration).Unix(),
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

func (ut *userToken) Validate(signedToken string) error {
	claims := userTokenClaims{}
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

	ut.uid = claims.UID

	return nil
}
