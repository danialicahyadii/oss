package request

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	LastName string `json:"last_name" binding:"max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=8"`
	IsActive *bool  `json:"is_active"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	LastName string `json:"last_name" binding:"max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
	IsActive *bool  `json:"is_active"`
}

type UserListRequest struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

type UpdateUserRolesRequest struct {
	RoleUUIDs []string `json:"role_uuids"`
}
