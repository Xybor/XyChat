package tokens

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xybor/xychat/helpers"
	xyerrors "github.com/xybor/xychat/xyerrors/v1"
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
	insecureSecret := "1n53cur37@k3n"
	secretString, err := helpers.ReadEnv("TOKEN_SECRET")
	if err != nil {
		secretString = insecureSecret
		log.Println("[Xychat] You are using the default TOKEN_SECRET. " +
			"Please set a secure value.")
	}

	secret := []byte(secretString)
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
func (ut userToken) Generate() (string, xyerrors.XyError) {
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
		log.Println(err)
		return "", xyerrors.ErrorUnknown
	}

	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Println(err)
		return "", xyerrors.ErrorCannotCreateToken.New("Something is wrong when creating token")
	}

	return signedToken, xyerrors.NoError
}

// Validate verifies the provided token and set id to the userToken.
func (ut *userToken) Validate(signedToken string) xyerrors.XyError {
	claims := userTokenClaims{}
	_, err := jwt.ParseWithClaims(
		signedToken,
		&claims,
		getSecret,
	)

	if err != nil {
		log.Println(err)
		return xyerrors.ErrorUnknown.New("Invalid token")
	}

	expiration := claims.StandardClaims.ExpiresAt
	if expiration-time.Now().Unix() < 0 {
		return xyerrors.ErrorInvalidToken.New("Expired token")
	}

	ut.id = claims.ID

	return xyerrors.NoError
}
