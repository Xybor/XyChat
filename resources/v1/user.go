package v1

type UserResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
	Age      *uint   `json:"age,omitempty"`
	Gender   *string `json:"gender,omitempty"`
}
type UserRegisterRequest struct {
	SrcId    *uint   `context:"id"`
	Username string  `json:"username" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Role     *string `json:"role"`
}

type UserAuthenticateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserProfileRequest struct {
	SrcID *uint `context:"id"`
}

type UserSelectRequest struct {
	SrcID *uint `context:"id"`
	DstID uint  `uri:"id" validate:"required"`
}

type UserUpdateRequest struct {
	SrcId  *uint   `context:"id"`
	DstID  uint    `uri:"id" validate:"required"`
	Age    *uint   `json:"age"`
	Gender *string `json:"gender"`
}

type UserUpdateRoleRequest struct {
	SrcId *uint  `context:"id"`
	DstID uint   `uri:"id" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type UserUpdatePasswordRequest struct {
	SrcId       *uint   `context:"id"`
	DstId       uint    `uri:"id" validate:"required"`
	OldPassword *string `json:"old_password"`
	NewPassword string  `json:"new_password" validate:"required"`
}
