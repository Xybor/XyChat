package v1

import "errors"

var (
	ErrorPermission                = errors.New("you don't have the permission")
	ErrorUnknownRole               = errors.New("unknown role")
	ErrorExistedUsername           = errors.New("username existed")
	ErrorInvalidOldPassword        = errors.New("invalid old password")
	ErrorFailedAuthentication      = errors.New("invalid username or password")
	ErrorDuplicatedConnection      = errors.New("duplicated connection")
	ErrorNotEnoughUserToCreateRoom = errors.New("not enough user to create room")
	ErrorCannotAccessToRoom        = errors.New("you cannot access to room")
	ErrorNotYetJoinInRoom          = errors.New("you had not yet joined in this room")
	ErrorNotInManagementList       = errors.New("service not in management list")
	ErrorUnknown                   = errors.New("some error occurs")
)
