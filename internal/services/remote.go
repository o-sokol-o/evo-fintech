package services

import (
	"context"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

type ServiceRemoteCSV struct {
	repo IRepoRemote
}

func NewRemoteServices(repo IRepoRemote) *ServiceRemoteCSV {
	return &ServiceRemoteCSV{repo: repo}
}

func (s *ServiceRemoteCSV) Get(ctx context.Context) ([]domain.Transaction, error) {
	return s.repo.Get(ctx)
}
