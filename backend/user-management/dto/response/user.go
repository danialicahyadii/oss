package response

import "time"

type UserResponse struct {
	UUID      string    `json:"uuid"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type UserRoleResponse struct {
	RoleUUID    string `json:"role_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
