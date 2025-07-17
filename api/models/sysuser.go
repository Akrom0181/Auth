package models

type SysUser struct {
	Id        string   `json:"id"`
	Status    string   `json:"status"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	CreatedAt string   `json:"created_at"`
	CreatedBy string   `json:"created_by"`
	Roles     []string `json:"roles"`
}

type GetSingleSysUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginSysUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetListSysUserRequest struct {
	Search string `json:"search"`
	Page   int64  `json:"page"`
	Limit  int64  `json:"limit"`
}

type GetListSysUserResponse struct {
	Items []SysUser `json:"items"`
	Count int64     `json:"count"`
}

type SysUserRoles struct {
	Id        string `json:"id"`
	SysUserId string `json:"sys_user_id"`
	RoleId    string `json:"role_id"`
}

type Roles struct {
	Id        string `json:"id"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}
