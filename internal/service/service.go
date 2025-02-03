//go:generate mockgen -destination=mock/service_mock.go -package=mock . Service
package service

import (
	"context"
	"time"

	"github.com/ferdiebergado/goweb/internal/repository"
)

type Service interface {
	PingDB(ctx context.Context) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) PingDB(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.PingContext(pingCtx)
}
