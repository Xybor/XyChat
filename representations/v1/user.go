package v1

type UserRepresentation struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
	Age      *uint   `json:"age,omitempty"`
	Gender   *string `json:"gender,omitempty"`
}
