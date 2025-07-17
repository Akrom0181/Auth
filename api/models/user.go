package models

type User struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type UserSingleRequest struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type GetListRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetListUserResponse struct {
	Items []User `json:"items"`
	Count int64  `json:"count"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}
