package response

type MeResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}
