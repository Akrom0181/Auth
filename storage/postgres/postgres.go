package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Akrom0181/Auth/config"
	"github.com/Akrom0181/Auth/pkg/logger"
	"github.com/Akrom0181/Auth/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Store struct {
	Pool   *pgxpool.Pool
	logger logger.ILogger
	cfg    config.Config
}

func New(ctx context.Context, cfg config.Config, logger logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(`host=%s port=%v user=%s password=%s database=%s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDatabase)

	pgPoolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pgPoolConfig.MaxConns = 100
	pgPoolConfig.MaxConnLifetime = time.Hour

	newPool, err := pgxpool.NewWithConfig(context.Background(), pgPoolConfig)
	if err != nil {
		fmt.Println("error while connecting to db", err.Error())
		return nil, err
	}

	return Store{
		Pool:   newPool,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (s Store) CloseDB() {
	s.Pool.Close()
}

func (s Store) User() storage.IUserStorage {
	newUser := NewUserRepo(s.Pool, s.logger)

	return &newUser
}

func (s Store) SysUser() storage.ISysUserStorage {
	newSysUser := NewSysUserRepo(s.Pool, s.logger, NewRoleRepo(s.Pool, s.logger))

	return &newSysUser
}

func (s Store) Otp() storage.IOtpStorage {
	newOtp := NewOtpRepo(s.Pool, s.logger)

	return &newOtp
}

func (s Store) Role() storage.IRoleStorage {
	newRole := NewRoleRepo(s.Pool, s.logger)

	return &newRole
}
