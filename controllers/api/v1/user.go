package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apihelpers "github.com/xybor/xychat/helpers/api"
	"github.com/xybor/xychat/helpers/context"
	"github.com/xybor/xychat/helpers/tokens"
	services "github.com/xybor/xychat/services/v1"
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
	if err.Errno() != 0 {
		role = services.RoleMember
	}

	userService := services.CreateUserService(id, true)

	xerr := userService.Register(username, password, role)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(nil)
	ctx.JSON(http.StatusCreated, response)
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
	userService := services.CreateUserService(nil, true)

	userRepresentation, xerr := userService.Authenticate(username, password)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Create a token with the expired duration is 24 hours
	userToken := tokens.CreateUserToken(userRepresentation.ID, 24*time.Hour)

	token, xerr := userToken.Generate()
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Set the token as a cookie.  It is configured as httponly for safety and
	// affects on all url.
	oneDay := 24 * 60 * 60
	context.SetCookie(ctx, "xytok", token, oneDay)

	response := apihelpers.NewAPIResponse(userRepresentation)
	ctx.JSON(http.StatusOK, response)
}

// UserProfileHandler handles an incoming request with one optional query
// parameter: token.
//
// Behavior: It gets the profile of the user determined by the token.
//
// Response: A profile + An error message.
func UserProfileHandler(ctx *gin.Context) {
	id := context.GetUID(ctx)

	userService := services.CreateUserService(id, true)

	profile, xerr := userService.SelfSelect()
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(profile)
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
	destId, xerr := context.GetURLParamAsUint(ctx, "id")
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(id, true)

	profile, xerr := userService.Select(destId)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(profile)
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

	age, xerr := context.RetrieveQueryAsPUint(ctx, "age")
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Gender has been already a string, so that it doesn't have any error to
	// handle.
	gender := context.RetrieveQueryAsPString(ctx, "gender")

	id := context.GetUID(ctx)
	destId, xerr := context.GetURLParamAsUint(ctx, "id")
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(id, true)
	xerr = userService.UpdateInfo(destId, age, gender)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(nil)
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

	destId, xerr := context.GetURLParamAsUint(ctx, "id")
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	role := context.MustRetrieveQuery(ctx, "role")

	userService := services.CreateUserService(id, true)
	xerr = userService.UpdateRole(destId, role)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(nil)
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

	destId, xerr := context.GetURLParamAsUint(ctx, "id")
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	oldpassword := context.RetrieveQueryAsPString(ctx, "oldpassword")
	newpassword := context.MustRetrieveQuery(ctx, "newpassword")

	userService := services.CreateUserService(id, true)
	xerr = userService.UpdatePassword(destId, oldpassword, newpassword)
	if xerr.Errno() != 0 {
		response := apihelpers.NewAPIError(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := apihelpers.NewAPIResponse(nil)
	ctx.JSON(http.StatusOK, response)
}
