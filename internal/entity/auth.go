package entity

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
}

type VerifyEmail struct {
	Email    string `json:"email"`
	Otp      string `json:"otp"`
	Platform string `json:"platform"`
}
