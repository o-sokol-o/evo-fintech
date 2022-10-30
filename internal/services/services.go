package services

import (
	"context"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

type IRepoEVO interface {
	GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error)
	InsertTransactions(ctx context.Context, transactions []domain.Transaction) error
}
type IRepoRemote interface {
	Get(ctx context.Context, from, to *int) ([]domain.Transaction, error)
}

type Services struct {
	ServicesEVO    *ServiceEVO
	ServicesRemote *ServiceRemoteCSV
}

func New(repoEVO IRepoEVO, repoRemote IRepoRemote) *Services {
	return &Services{
		ServicesEVO:    NewEvoServices(repoEVO),
		ServicesRemote: NewRemoteServices(repoRemote),
	}
}
