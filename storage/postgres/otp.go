package postgres

import (
	"context"

	"github.com/Akrom0181/Auth/api/models"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OtpRepo struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewOtpRepo(db *pgxpool.Pool, logger logger.ILogger) OtpRepo {
	return OtpRepo{
		db:     db,
		logger: logger,
	}
}

func (r *OtpRepo) Create(ctx context.Context, req models.Otp) (models.Otp, error) {

	query := `INSERT INTO otp (id, email, status, code, expires_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, req.Id, req.Email, req.Status, req.Code, req.ExpiresAt)
	if err != nil {
		r.logger.Error("failed to create otp", logger.Error(err))
		return models.Otp{}, err
	}

	return req, nil
}

func (r *OtpRepo) GetSingle(ctx context.Context, req models.GetSingleOTP) (models.Otp, error) {
	query := `SELECT id, email, status, code, expires_at FROM otp WHERE id = $1 AND (status = $2 OR email = $3)`
	row := r.db.QueryRow(ctx, query, req.Id, req.Status, req.Email)
	var otp models.Otp
	err := row.Scan(&otp.Id, &otp.Email, &otp.Status, &otp.Code, &otp.ExpiresAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.logger.Info("otp not found", logger.String("id", req.Id))
			return models.Otp{}, nil
		}
		r.logger.Error("failed to get otp", logger.Error(err))
		return models.Otp{}, err
	}
	return otp, nil
}

func (r *OtpRepo) Update(ctx context.Context, req models.Otp) error {
	query := `UPDATE otp SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, req.Status, req.Id)
	if err != nil {
		r.logger.Error("failed to update otp", logger.Error(err))
		return err
	}
	return nil
}
