package service

import (
	"context"

	"github.com/ferdiebergado/goweb/internal/repository"
)

type Service struct {
	repo repository.DBRepo
}

func NewDBService(repo repository.DBRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) DBVersion(ctx context.Context) (string, error) {
	return s.repo.DBVersion(ctx)
}
