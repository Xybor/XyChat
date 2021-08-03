package v1

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ctrl "github.com/xybor/xychat/controllers"
	"github.com/xybor/xychat/helpers"
	apihelper "github.com/xybor/xychat/helpers/api/v1"
	"github.com/xybor/xychat/helpers/tokens"
	service "github.com/xybor/xychat/services/v1"
)

func RegisterUserHandler(c *gin.Context) {
	ctrl.SetReceivingMethod(ctrl.GET)

	username, err := ctrl.GetParam(c, "username")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.LackOfInput, "empty username")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	password, err := ctrl.GetParam(c, "password")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.LackOfInput, "empty password")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userService := service.CreateUserService(nil, &username, &password)

	err = userService.Register()
	if err != nil {
		log.Println(err)
		response := apihelper.NewAPIError(ctrl.FailedProcess, "registration failed")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewEmptyAPIResponse()
	c.JSON(http.StatusAccepted, response)
}

func AuthenticateUserHandler(c *gin.Context) {
	ctrl.SetReceivingMethod(ctrl.GET)

	username, err := ctrl.GetParam(c, "username")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.LackOfInput, "empty username")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	password, err := ctrl.GetParam(c, "password")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.LackOfInput, "empty password")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userService := service.CreateUserService(nil, &username, &password)

	userRepresentation, err := userService.Authenticate()
	if err != nil {
		log.Println(err)
		response := apihelper.NewAPIError(ctrl.FailedProcess, "authentication failed")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userToken := tokens.CreateUserToken(userRepresentation.ID, 15*time.Minute)

	token, err := userToken.Generate()
	if err != nil {
		log.Println(err)
		response := apihelper.NewAPIError(ctrl.FailedProcess, "couldn't create token")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.SetCookie(
		"auth",                                 // name
		token,                                  // value
		60*60*24*7,                             // maxage (7 days)
		"/",                                    // path
		helpers.ReadEnv("domain", "localhost"), // domain
		false,                                  // secure
		true,                                   // httponly
	)

	userRepresentation.Token = &token
	response := apihelper.NewAPIResponse(userRepresentation)
	c.JSON(http.StatusAccepted, response)
}

func GetProfileHandler(c *gin.Context) {
	suid, ok := c.Get("UID")

	if !ok {
		response := apihelper.NewAPIError(ctrl.Unauthenticated, "Unauthenticated")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	uid := suid.(uint)

	userService := service.CreateUserService(&uid, nil, nil)

	profile, err := userService.GetProfile()
	if err != nil {
		log.Println(err)
		response := apihelper.NewAPIError(ctrl.FailedProcess, "Couldn't get user profile")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewAPIResponse(profile)
	c.JSON(http.StatusOK, response)
}
