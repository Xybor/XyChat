package v1

import "errors"

var (
	ErrorPermission           = errors.New("you don't have the permission")
	ErrorUnknownRole          = errors.New("unknown role")
	ErrorExistedUsername      = errors.New("username existed")
	ErrorInvalidOldPassword   = errors.New("invalid old password")
	ErrorFailedAuthentication = errors.New("invalid username or password")
	ErrorDuplicatedConnection = errors.New("duplicated connection")
	ErrorUnknown              = errors.New("some error occurs")
)
