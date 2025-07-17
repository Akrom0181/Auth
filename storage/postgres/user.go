package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewUserRepo(db *pgxpool.Pool, logger logger.ILogger) UserRepo {
	return UserRepo{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req models.User) (models.User, error) {
	req.Id = uuid.NewString()

	query := `INSERT INTO users (id, status, email, password, name, created_at) VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := r.db.Exec(ctx, query, req.Id, req.Status, req.Email, req.Password, req.Name)
	if err != nil {
		r.logger.Error("failed to create user", logger.Error(err))
		return models.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req models.UserSingleRequest) (models.User, error) {
	var response models.User

	var (
		query     string
		args      []interface{}
		createdAt time.Time
	)

	switch {
	case req.Id != "":
		query = `SELECT id, email, password, name, status, created_at FROM users WHERE id = $1`
		args = append(args, req.Id)
	case req.Email != "" && req.Status != "":
		query = `SELECT id, email, password, name, status, created_at FROM users WHERE email = $1 AND status = $2`
		args = append(args, req.Email, req.Status)
	case req.Email != "":
		query = `SELECT id, email, password, name, status, created_at FROM users WHERE email = $1`
		args = append(args, req.Email)
	default:
		return models.User{}, fmt.Errorf("invalid request: must provide id or (email and status) or email")
	}

	row := r.db.QueryRow(ctx, query, args...)
	err := row.Scan(&response.Id, &response.Email, &response.Password, &response.Name, &response.Status, &createdAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Info("user not found", logger.String("id", req.Id), logger.String("email", req.Email), logger.String("status", req.Status))
			return models.User{}, nil
		}
		r.logger.Error("failed to get user", logger.Error(err))
		return models.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)

	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context, req models.GetListRequest) (models.GetListUserResponse, error) {

	var (
		response  models.GetListUserResponse
		createdAt time.Time
	)

	query := `SELECT id, email, name, created_at FROM users WHERE email ILIKE '%' || $1 || '%' LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, req.Search, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		r.logger.Error("failed to get user list", logger.Error(err))
		return response, err
	}

	defer rows.Close()
	for rows.Next() {
		var item models.User
		err := rows.Scan(&item.Id, &item.Email, &item.Name, &item.CreatedAt)
		if err != nil {
			r.logger.Error("failed to scan user row", logger.Error(err))
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery := `SELECT COUNT(*) FROM users WHERE email ILIKE '%' || $1 || '%'`
	err = r.db.QueryRow(ctx, countQuery, req.Search).Scan(&response.Count)
	if err != nil {
		r.logger.Error("failed to count users", logger.Error(err))
		return response, err
	}

	return response, nil
}

func (r *UserRepo) Update(ctx context.Context, req models.User) (models.User, error) {
	mp := map[string]interface{}{
		"email": req.Email,
		"name":  req.Name,
	}

	if req.Password != "" {
		mp["password"] = req.Password
	}

	setClause := ""
	args := []interface{}{}
	argPos := 1

	for key, val := range mp {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, val)
		argPos++
	}

	args = append(args, req.Id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", setClause, argPos)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to update user", logger.Error(err))
		return models.User{}, err
	}

	return req, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user", logger.Error(err))
		return err
	}
	if err == pgx.ErrNoRows {
		r.logger.Info("user not found for deletion", logger.String("id", id))
		return err
	}
	return nil
}
