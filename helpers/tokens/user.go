package tokens

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xybor/xychat/helpers"
)

var (
	ErrorExpiredToken = errors.New("token is expired")
	ErrorInvalidToken = errors.New("token is invalid")
)

type userToken struct {
	id         uint
	expiration time.Duration
}

type userTokenClaims struct {
	ID uint
	*jwt.StandardClaims
}

// getSecret reads token_secret from environment variables.
func getSecret(token *jwt.Token) (interface{}, error) {
	secret_string := helpers.ReadEnv("token_secret", "1ns3cur3t0k3n")

	secret := []byte(secret_string)
	return secret, nil
}

// CreateUserToken creates a userToken with provided information.
func CreateUserToken(id uint, expiration time.Duration) userToken {
	return userToken{
		id:         id,
		expiration: expiration,
	}
}

// CreateEmptyUserToken creates an empty userToken for the Validate() method.
func CreateEmptyUserToken() userToken {
	return userToken{}
}

// GetUID return id in the userToken.
func (ut userToken) GetUID() uint {
	return ut.id
}

// Generate creates a token with the userToken information.
func (ut userToken) Generate() (string, error) {
	claims := userTokenClaims{
		ID: ut.id,
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

// Validate verifies the provided token and set id to the userToken.
func (ut *userToken) Validate(signedToken string) error {
	claims := userTokenClaims{}
	_, err := jwt.ParseWithClaims(
		signedToken,
		&claims,
		getSecret,
	)

	if err != nil {
		log.Print(err)
		return ErrorInvalidToken
	}

	expiration := claims.StandardClaims.ExpiresAt
	if expiration-time.Now().Unix() < 0 {
		return ErrorExpiredToken
	}

	ut.id = claims.ID

	return nil
}
