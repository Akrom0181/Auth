package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewRoleRepo(db *pgxpool.Pool, logger logger.ILogger) RoleRepo {
	return RoleRepo{
		db:     db,
		logger: logger,
	}
}

func (r *RoleRepo) Create(ctx context.Context, req models.Role) (models.Role, error) {
	req.Id = uuid.NewString()

	query := `INSERT INTO roles (id, name, status, created_by, created_at) VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.db.Exec(ctx, query, req.Id, req.Name, req.Status, req.CreatedBy)
	if err != nil {
		r.logger.Error("failed to create role", logger.Error(err))
		return models.Role{}, err
	}

	return req, nil
}

func (r *RoleRepo) GetSingle(ctx context.Context, req models.ID) (models.Role, error) {
	var createdAt time.Time

	query := `SELECT id, name, status, created_by, created_at FROM roles WHERE id = $1`
	row := r.db.QueryRow(ctx, query, req.Id)
	var role models.Role
	err := row.Scan(&role.Id, &role.Name, &role.Status, &role.CreatedBy, &createdAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Info("role not found", logger.String("id", req.Id))
			return models.Role{}, nil
		}
		r.logger.Error("failed to get role", logger.Error(err))
		return models.Role{}, err

	}
	role.CreatedAt = createdAt.Format(time.RFC3339)
	return role, nil
}

func (r *RoleRepo) GetList(ctx context.Context, req models.GetListRequest) (models.GetListRoleResponse, error) {
	var response = models.GetListRoleResponse{}

	query := `SELECT id, name, status, created_by, created_at FROM roles WHERE name ILIKE '%' || $1 || '%' LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, req.Search, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		r.logger.Error("failed to get roles list", logger.Error(err))
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item      models.Role
			createdAt time.Time
			createdBy sql.NullString
		)

		err := rows.Scan(&item.Id, &item.Name, &item.Status, &createdBy, &createdAt)
		if err != nil {
			r.logger.Error("failed to scan role", logger.Error(err))
			return response, err
		}

		// Agar created_by NULL bo‘lsa, bo‘sh string beriladi
		if createdBy.Valid {
			item.CreatedBy = createdBy.String
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		response.Items = append(response.Items, item)
	}

	countQuery := `SELECT COUNT(*) FROM roles WHERE name ILIKE '%' || $1 || '%'`
	err = r.db.QueryRow(ctx, countQuery, req.Search).Scan(&response.Count)
	if err != nil {
		r.logger.Error("failed to count roles", logger.Error(err))
		return response, err
	}

	return response, nil
}

func (r *RoleRepo) Update(ctx context.Context, req models.Role) (models.Role, error) {
	query := `UPDATE roles SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, req.Name, req.Id)
	if err != nil {
		r.logger.Error("failed to update role", logger.Error(err))
		return models.Role{}, err
	}

	return req, nil
}

func (r *RoleRepo) Delete(ctx context.Context, req models.ID) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.Exec(ctx, query, req.Id)
	if err != nil {
		r.logger.Error("failed to delete role", logger.Error(err))
		return err
	}
	return nil
}

func (r *RoleRepo) ExistsByIDAndStatus(ctx context.Context, id string, status string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1 AND status = $2)`
	err := r.db.QueryRow(ctx, query, id, status).Scan(&exists)
	if err != nil {
		r.logger.Error("failed to check role existence", logger.Error(err))
		return false, err
	}
	return exists, nil
}
