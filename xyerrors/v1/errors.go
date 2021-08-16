package v1

import "net/http"

var (
	NoError = NewError(http.StatusOK, 0)

	ErrorUnknown = NewError(http.StatusInternalServerError, 10001).New(
		"Something's wrong, please contact the administrator or mod")

	ErrorUnknownInput              = NewError(http.StatusUnprocessableEntity, 10002)
	ErrorNotAuthenticated          = NewError(http.StatusUnauthorized, 10003)
	ErrorPermission                = NewError(http.StatusUnauthorized, 10004)
	ErrorExistedUsername           = NewError(http.StatusConflict, 10005)
	ErrorFailedAuthentication      = NewError(http.StatusUnauthorized, 10006)
	ErrorDuplicatedConnection      = NewError(http.StatusConflict, 10007)
	ErrorNotEnoughUserToCreateRoom = NewError(http.StatusUnprocessableEntity, 10008)
	ErrorNotYetJoinInRoom          = NewError(http.StatusUnprocessableEntity, 10009)
	ErrorNotInManagementList       = NewError(http.StatusInternalServerError, 10010)
	ErrorInvalidService            = NewError(http.StatusInternalServerError, 10012)
	ErrorNotFoundInput             = NewError(http.StatusUnprocessableEntity, 10013)
	ErrorSyntaxInput               = NewError(http.StatusBadRequest, 10014)
	ErrorMediaType                 = NewError(http.StatusUnsupportedMediaType, 10015)
	ErrorCannotCreateToken         = NewError(http.StatusInternalServerError, 10016)
	ErrorInvalidToken              = NewError(http.StatusUnprocessableEntity, 10017)
	ErrorCannotUpgradeToWebsocket  = NewError(http.StatusMethodNotAllowed, 10018)
)
