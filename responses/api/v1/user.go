package v1

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Token    *string `json:"token,omitempty"`
}
