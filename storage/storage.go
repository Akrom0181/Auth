package storage

import (
	"context"

	"github.com/Akrom0181/Auth/api/models"
)

type IStorage interface {
	CloseDB()
	User() IUserStorage
	SysUser() ISysUserStorage
	Otp() IOtpStorage
	Role() IRoleStorage
}

type (
	// UserStorage -.

	IUserStorage interface {
		Create(ctx context.Context, req models.User) (models.User, error)
		GetSingle(ctx context.Context, req models.UserSingleRequest) (models.User, error)
		GetList(ctx context.Context, req models.GetListRequest) (models.GetListUserResponse, error)
		Update(ctx context.Context, req models.User) (models.User, error)
		Delete(ctx context.Context, id string) error
	}

	// SysUserStorage -.

	ISysUserStorage interface {
		Create(ctx context.Context, user models.SysUser) error
		AttachRole(ctx context.Context, userID, roleID string) error
		GetByEmailAndStatus(ctx context.Context, email string, statuses []string) (models.SysUser, error)
		GetSingle(ctx context.Context, req models.GetSingleSysUser) (models.SysUser, error)
	}

	// OtpStorage -.

	IOtpStorage interface {
		Create(ctx context.Context, req models.Otp) (models.Otp, error)
		GetSingle(ctx context.Context, req models.GetSingleOTP) (models.Otp, error)
		Update(ctx context.Context, req models.Otp) error
	}

	// RoleStorage -.

	IRoleStorage interface {
		Create(ctx context.Context, req models.Role) (models.Role, error)
		GetSingle(ctx context.Context, req models.ID) (models.Role, error)
		GetList(ctx context.Context, req models.GetListRequest) (models.GetListRoleResponse, error)
		Update(ctx context.Context, req models.Role) (models.Role, error)
		Delete(ctx context.Context, id models.ID) error
		ExistsByIDAndStatus(ctx context.Context, id string, status string) (bool, error)
	}
)
