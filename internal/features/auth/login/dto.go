package login

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required,len=6"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	UserID      int64    `json:"user_id"`
	FullName    string   `json:"full_name"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
}
