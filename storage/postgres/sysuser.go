package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SysUserRepo struct {
	db       *pgxpool.Pool
	logger   logger.ILogger
	roleRepo RoleRepo
}

func NewSysUserRepo(db *pgxpool.Pool, logger logger.ILogger, roleRepo RoleRepo) SysUserRepo {
	return SysUserRepo{
		db:       db,
		logger:   logger,
		roleRepo: roleRepo,
	}
}

func (r *SysUserRepo) Create(ctx context.Context, user models.SysUser) error {
	query := `
		INSERT INTO sysusers (id, name, email, password, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	_, err := r.db.Exec(ctx, query, user.Id, user.Name, user.Email, user.Password, user.Status)
	if err != nil {
		r.logger.Error("failed to insert sysuser", logger.Error(err))
	}
	return err
}

func (r *SysUserRepo) AttachRole(ctx context.Context, userID, roleID string) error {
	id := uuid.NewString()
	query := `
		INSERT INTO sysuser_roles (id, sysuser_id, role_id, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := r.db.Exec(ctx, query, id, userID, roleID)
	if err != nil {
		r.logger.Error("failed to attach role to sysuser", logger.Error(err))
	}
	return err
}

func (r *SysUserRepo) GetByEmailAndStatus(ctx context.Context, email string, statuses []string) (models.SysUser, error) {
	user := models.SysUser{}
	query := `SELECT id, name, email, password, status FROM sysusers WHERE email = $1 AND status = ANY($2) LIMIT 1`
	row := r.db.QueryRow(ctx, query, email, statuses)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("sysuser not found by email", logger.String("email", email))
			return models.SysUser{}, errors.New("user not found")
		}
		r.logger.Error("failed to get sysuser by email", logger.Error(err))
		return models.SysUser{}, err
	}
	return user, nil
}

func (r *SysUserRepo) GetSingle(ctx context.Context, req models.GetSingleSysUser) (models.SysUser, error) {
	user := models.SysUser{}
	var query string
	var args []interface{}
	if req.Id != "" {
		query = `SELECT id, name, email, password, status FROM sysusers WHERE id = $1`
		args = append(args, req.Id)
	} else if req.Email != "" {
		query = `SELECT id, name, email, password, status FROM sysusers WHERE email = $1`
		args = append(args, req.Email)
	} else {
		return models.SysUser{}, errors.New("invalid request: must provide id or (email and status) or email")
	}
	row := r.db.QueryRow(ctx, query, args...)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.SysUser{}, nil
		}
		r.logger.Error("failed to get sysuser", logger.Error(err))
		return models.SysUser{}, err
	}
	return user, nil
}

func (r *SysUserRepo) CreateSuperAdmin(ctx context.Context) error {
	const superAdminEmail = "super@admin.com"
	const superAdminPassword = "SuperSecurePassword123"
	const superAdminRoleID = "00000000-0000-0000-0000-000000000001"

	exists, err := r.roleRepo.ExistsByIDAndStatus(ctx, superAdminRoleID, "active")
	if err != nil {
		return err
	}
	if !exists {
		_, err := r.roleRepo.Create(ctx, models.Role{
			Id:     superAdminRoleID,
			Name:   "super_admin",
			Status: "active",
		})
		if err != nil {
			return err
		}
		log.Println("Created super_admin role")
	}

	user, err := r.GetByEmailAndStatus(ctx, superAdminEmail, []string{"active", "blocked"})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if user.Id != "" {
		log.Println("Super admin already exists.")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.SysUser{
		Id:       uuid.New().String(),
		Name:     "Super Admin",
		Email:    superAdminEmail,
		Password: string(hashedPassword),
		Status:   "active",
	}
	if err := r.Create(ctx, admin); err != nil {
		return err
	}
	log.Println("Created super_admin user")

	if err := r.AttachRole(ctx, admin.Id, superAdminRoleID); err != nil {
		return err
	}
	log.Println("Assigned super_admin role")

	return nil
}
