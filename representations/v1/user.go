package v1

type UserRepresentation struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
	Age      *int    `json:"age,omitempty"`
	Gender   *string `json:"gender,omitempty"`
	Token    *string `json:"token,omitempty"`
}
