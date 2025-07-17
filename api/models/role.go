package models

type Role struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type GetListRoleRequest struct {
	Search string `json:"search"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

type GetListRoleResponse struct {
	Items []Role `json:"items"`
	Count int64  `json:"count"`
}
