package v1

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ctr "github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers"
	apihelpersv1 "github.com/xybor/xychat/helpers/api/v1"
	"github.com/xybor/xychat/helpers/tokens"
	sv "github.com/xybor/xychat/services/api/v1"
)

func RegisterUserHandler(c *gin.Context) {
	ctr.SetReceivingMethod(ctr.GET)

	username, err := ctr.GetParam(c, "username")
	if err != nil {
		response := apihelpersv1.NewAPIError(LackOfInput, "empty username")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	password, err := ctr.GetParam(c, "password")
	if err != nil {
		response := apihelpersv1.NewAPIError(LackOfInput, "empty password")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userService := sv.UserService{
		Username: &username,
		Password: &password,
	}

	err = userService.RegisterUser()
	if err != nil {
		log.Println(err)
		response := apihelpersv1.NewAPIError(FailedProcess, "registration failed")
		c.JSON(http.StatusNotAcceptable, response)
		return
	}

	response := apihelpersv1.NewAPIResponse(nil)
	c.JSON(http.StatusAccepted, response)
}

func AuthenticateUserHandler(c *gin.Context) {
	ctr.SetReceivingMethod(ctr.GET)

	username, err := ctr.GetParam(c, "username")
	if err != nil {
		response := apihelpersv1.NewAPIError(LackOfInput, "empty username")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	password, err := ctr.GetParam(c, "password")
	if err != nil {
		response := apihelpersv1.NewAPIError(LackOfInput, "empty password")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userService := sv.UserService{
		Username: &username,
		Password: &password,
	}

	userResponse, err := userService.Authenticate()
	if err != nil {
		log.Println(err)
		response := apihelpersv1.NewAPIError(FailedProcess, "authentication failed")
		c.JSON(http.StatusNotAcceptable, response)
		return
	}

	userToken := tokens.UserToken{
		ID:         userResponse.ID,
		Expiration: 15 * time.Minute,
	}

	token, err := userToken.Create()
	if err != nil {
		log.Println(err)
		response := apihelpersv1.NewAPIError(FailedProcess, "couldn't create token")
		c.JSON(http.StatusNotAcceptable, response)
		return
	}

	userResponse.Token = &token

	response := apihelpersv1.NewAPIResponse(userResponse)
	c.SetCookie(
		"auth",                                 // name
		token,                                  // value
		60*60*24*7,                             // maxage (7 days)
		"/",                                    // path
		helpers.ReadEnv("domain", "localhost"), // domain
		false,                                  // secure
		true,                                   // httponly
	)
	c.JSON(http.StatusAccepted, response)
}

func GetProfileHandler(c *gin.Context) {
	suid, ok := c.Get("UID")

	if !ok {
		response := apihelpersv1.NewAPIError(Unauthenticated, "Unauthenticated")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	uid := suid.(uint)

	us := sv.UserService{
		ID: &uid,
	}

	profile, err := us.GetProfile()
	if err != nil {
		log.Println(err)
		response := apihelpersv1.NewAPIError(FailedProcess, "Couldn't get user profile")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelpersv1.NewAPIResponse(profile)
	c.JSON(http.StatusOK, response)
}
