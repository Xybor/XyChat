package v1

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ctrl "github.com/xybor/xychat/controllers"
	apihelper "github.com/xybor/xychat/helpers/api/v1"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/helpers/tokens"
	service "github.com/xybor/xychat/services/v1"
)

// UserRegisterHandler handles an incoming request with four query parameters:
// username, password, and (optional) token, role.
//
// Behavior: It registers a user account with the provided information.  If the
// role is a member, it doesn't need a subject to register; otherwise, the
// subject must be higher role.
//
// Response: An error message.
func UserRegisterHandler(ctx *gin.Context) {
	context.SetRetrievingMethod(context.POST)

	id := context.GetUID(ctx)
	username := context.MustRetrieveQuery(ctx, "username")
	password := context.MustRetrieveQuery(ctx, "password")

	// By default, role is set as 'member' if it isn't provided.
	role, err := context.RetrieveQuery(ctx, "role")
	if err != nil {
		role = service.RoleMember
	}

	userService := service.CreateUserService(id)

	err = userService.Register(username, password, role)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewEmptyAPIResponse()
	ctx.JSON(http.StatusAccepted, response)
}

// UserAuthenticateHandler handles an incoming request with three query
// parameters: username, password, and (optional) token.
//
// Behavior: It authenticates the user account with the provided information
// and generates an authenticated token.
//
// Response: A token in the body (for debugging) and as a cookie + An error
// message.
func UserAuthenticateHandler(ctx *gin.Context) {
	context.SetRetrievingMethod(context.POST)

	username := context.MustRetrieveQuery(ctx, "username")
	password := context.MustRetrieveQuery(ctx, "password")

	// Authentication doesn't need a subject to call the service
	userService := service.CreateUserService(nil)

	userRepresentation, err := userService.Authenticate(username, password)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// Create a token with the expired duration is 24 hours
	userToken := tokens.CreateUserToken(userRepresentation.ID, 5*time.Second)

	token, err := userToken.Generate()
	if err != nil {
		log.Println(err)
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, "couldn't create token")
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Set the token as a cookie.  It is configured as httponly for safety and
	// affects on all url.
	dayTime := 24 * 60 * 60
	context.SetCookie(ctx, "xytok", token, dayTime)

	response := apihelper.NewEmptyAPIResponse()
	ctx.JSON(http.StatusAccepted, response)
}

// UserProfileHandler handles an incoming request with one optional query
// parameter: token.
//
// Behavior: It gets the profile of the user determined by the token.
//
// Response: A profile + An error message.
func UserProfileHandler(ctx *gin.Context) {
	id := context.GetUID(ctx)

	userService := service.CreateUserService(id)

	profile, err := userService.SelfSelect()
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewAPIResponse(profile)
	ctx.JSON(http.StatusOK, response)
}

// UserGETHandler handles an incoming GET request of url /users/:id with one
// optional query parameter: token.
//
// Behavior: It gets the profile of a user determined by URL parameter id.  The
// subject is determined by the token.
//
// Response: A profile + An error message.
func UserGETHandler(ctx *gin.Context) {
	id := context.GetUID(ctx)
	destId, err := context.GetURLParamAsUint(ctx, "id")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorInput, "couldn't parse id")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userService := service.CreateUserService(id)

	profile, err := userService.Select(destId)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewAPIResponse(profile)
	ctx.JSON(http.StatusOK, response)
}

// UserPUTHandler handles an incoming PUT request of url /users/:id with three
// optional query parameters: token, age, gender.
//
// Behavior: It updates age and gender of the current profile.  The subject is
// determined by token.  The affected object is determined by the URL
// parameter.  If the request doesn't provide age or gender, they will be set
// to null (a.k.a be deleted).
//
// Response: An error message.
func UserPUTHandler(ctx *gin.Context) {
	context.SetRetrievingMethod(context.POST)

	age, err := context.RetrieveQueryAsPUint(ctx, "age")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorInvalidInput, "invalid age")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Gender has been already a string, so that it doesn't have any error to
	// handle.
	gender := context.RetrieveQueryAsPString(ctx, "gender")

	id := context.GetUID(ctx)
	destId, err := context.GetURLParamAsUint(ctx, "id")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorInput, "couldn't parse id")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userService := service.CreateUserService(id)
	err = userService.UpdateInfo(destId, age, gender)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewEmptyAPIResponse()
	ctx.JSON(http.StatusOK, response)
}

// UserChangRoleHandler handles an incoming PUT request of url /users/:id/role
// with two query parameters: role and (optional) token.
//
// Behavior: It updates role of the current profile.  The subject is determined
// by token.  The affected object is determined by the URL parameter.
//
// Response: An error message.
func UserChangeRoleHandler(ctx *gin.Context) {
	context.SetRetrievingMethod(context.POST)

	id := context.GetUID(ctx)

	destId, err := context.GetURLParamAsUint(ctx, "id")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorInput, "couldn't parse id")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	role := context.MustRetrieveQuery(ctx, "role")

	userService := service.CreateUserService(id)
	err = userService.UpdateRole(destId, role)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewEmptyAPIResponse()
	ctx.JSON(http.StatusOK, response)
}

// UserChangPasswordHandler handles an incoming PUT request of url
// /users/:id/password with three query parameters: newpassword and (optional)
// token, oldpassword.
//
// Behavior: It updates password of the current profile.  The subject is
// determined by token.  The affected object is determined by the URL
// parameter.  If the subject is higher role than the object, it doesn't need
// the oldpassword.
//
// Response: An error message.
func UserChangePasswordHandler(ctx *gin.Context) {
	context.SetRetrievingMethod(context.POST)

	id := context.GetUID(ctx)

	destId, err := context.GetURLParamAsUint(ctx, "id")
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorInput, "couldn't parse id")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	oldpassword := context.RetrieveQueryAsPString(ctx, "oldpassword")
	newpassword := context.MustRetrieveQuery(ctx, "newpassword")

	userService := service.CreateUserService(id)
	err = userService.UpdatePassword(destId, oldpassword, newpassword)
	if err != nil {
		response := apihelper.NewAPIError(ctrl.ErrorFailedProcess, err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := apihelper.NewEmptyAPIResponse()
	ctx.JSON(http.StatusOK, response)
}
