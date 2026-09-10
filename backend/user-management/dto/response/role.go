package response

type RoleResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type RolePermissionResponse struct {
	PermissionUUID string  `json:"permission_uuid"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
}
