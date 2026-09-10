package request

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// type ChangePasswordRequest struct {
// 	OldPassword     string `json:"old_password" binding:"required"`
// 	NewPassword     string `json:"new_password" binding:"required,min=8"`
// 	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
// }

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	DeviceUUID *string `json:"device_uuid,omitempty"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
