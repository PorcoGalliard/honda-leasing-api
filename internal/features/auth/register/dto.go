package register

var AllowedRegisterRoles = map[string]bool{
	"CUSTOMER":     true,
	"SALES":        true,
	"SURVEYOR":     true,
	"FINANCE":      true,
	"COLLECTION":   true,
	"ADMIN_CABANG": true,
}

type RegisterRequest struct {
	FullName    string `json:"full_name" binding:"required,min=3,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"required,min=8"`
	PIN         string `json:"pin" binding:"required,len=6"`
	Role        string `json:"role" binding:"omitempty"` // Opsional, default: CUSTOMER
}

type RegisterResponse struct {
	UserID      int64  `json:"user_id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
