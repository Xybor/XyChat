package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xybor/xychat/helpers"
	"github.com/xybor/xychat/helpers/tokens"
	"github.com/xybor/xychat/helpers/xybinders"
	resources "github.com/xybor/xychat/resources/v1"
	services "github.com/xybor/xychat/services/v1"
)

// UserRegisterHandler handles an incoming request and registers a user if it
// is valid.
func UserRegisterHandler(ctx *gin.Context) {
	request := new(resources.UserRegisterRequest)
	xerr := xybinders.Bind(ctx, request, xybinders.Json, xybinders.Context)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// By default, role is set as 'member' if it isn't provided.
	if request.Role == nil {
		role := services.RoleMember
		request.Role = &role
	}

	userService := services.CreateUserService(request.SrcId, true)

	xerr = userService.Register(request.Username, request.Password, *request.Role)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(nil)
	ctx.JSON(http.StatusCreated, response)
}

// UserAuthenticateHandler handles an incoming request and responds a cookie
// containing authenticated token if it is valid.
func UserAuthenticateHandler(ctx *gin.Context) {
	request := new(resources.UserAuthenticateRequest)
	xerr := xybinders.Bind(ctx, request, xybinders.Json)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Authentication doesn't need a subject to call the service
	userService := services.CreateUserService(nil, true)

	userRepresentation, xerr := userService.Authenticate(request.Username, request.Password)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Create a token with the expired duration is 24 hours
	userToken := tokens.CreateUserToken(userRepresentation.ID, 24*time.Hour)

	token, xerr := userToken.Generate()
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	// Set the token as a cookie.  It is configured as httponly for safety and
	// affects on all url.
	oneDay := 24 * 60 * 60
	helpers.SetCookie(ctx, "xytok", token, oneDay)

	response := helpers.NewResponse(userRepresentation)
	ctx.JSON(http.StatusOK, response)
}

// UserProfileHandler handles an incoming request and responds the current
// user's profile.
func UserProfileHandler(ctx *gin.Context) {
	request := resources.UserProfileRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(request.SrcID, true)

	profile, xerr := userService.SelfSelect()
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(profile)
	ctx.JSON(http.StatusOK, response)
}

// UserSelectHandler handles an incoming request and responds the profile of a
// specific user.
func UserSelectHandler(ctx *gin.Context) {
	request := resources.UserSelectRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context, xybinders.Uri)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(request.SrcID, true)

	profile, xerr := userService.Select(request.DstID)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(profile)
	ctx.JSON(http.StatusOK, response)
}

// UserUpdateHandler handles an incoming request and updates age and gender of
// a specific user.
func UserUpdateHandler(ctx *gin.Context) {
	request := resources.UserUpdateRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context, xybinders.Uri, xybinders.Json)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(request.SrcId, true)
	xerr = userService.UpdateInfo(request.DstID, request.Age, request.Gender)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(nil)
	ctx.JSON(http.StatusOK, response)
}

// UserChangRoleHandler handles an incoming request and updates role of a
// specific user.
func UserUpdateRoleHandler(ctx *gin.Context) {
	request := resources.UserUpdateRoleRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context, xybinders.Json, xybinders.Uri)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(request.SrcId, true)
	xerr = userService.UpdateRole(request.DstID, request.Role)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(nil)
	ctx.JSON(http.StatusOK, response)
}

// UserChangPasswordHandler handles an incoming request and changes the
// password of a specific user.
func UserChangePasswordHandler(ctx *gin.Context) {
	request := resources.UserUpdatePasswordRequest{}
	xerr := xybinders.Bind(ctx, &request, xybinders.Context, xybinders.Uri, xybinders.Json)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	userService := services.CreateUserService(request.SrcId, true)
	xerr = userService.UpdatePassword(request.DstId, request.OldPassword, request.NewPassword)
	if xerr.Errno() != 0 {
		response := helpers.NewErrorResponse(xerr)
		ctx.JSON(xerr.StatusCode(), response)
		return
	}

	response := helpers.NewResponse(nil)
	ctx.JSON(http.StatusOK, response)
}
