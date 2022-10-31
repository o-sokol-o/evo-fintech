package services

import (
	"context"
	"sync"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

type ServiceEVO struct {
	repo           IRepoEVO
	downloadStatus sync.Map
}

func NewEvoServices(repo IRepoEVO) *ServiceEVO {
	se := ServiceEVO{repo: repo}
	se.downloadStatus.Store("downloadStatus", domain.Unknown)
	return &se
}

func (s *ServiceEVO) GetFilteredData(ctx context.Context, input domain.FilterSearchInput) ([]domain.Transaction, error) {
	return s.repo.GetFilteredData(ctx, input)
}

func (s *ServiceEVO) FetchExternTransactions(ctx context.Context, url string) (domain.Status, error) {
	downloadStatus := domain.Unknown
	if ds, ok := s.downloadStatus.Load("downloadStatus"); ok {
		downloadStatus = ds.(domain.Status)
	}

	if url == "" {
		return downloadStatus, nil
	}

	if downloadStatus == domain.Processing || downloadStatus == domain.Skip {
		if downloadStatus == domain.Processing {
			s.downloadStatus.Store("downloadStatus", domain.Skip)
			return domain.Skip, nil
		}
		return downloadStatus, nil
	}

	// We request a list of transactions from an external service via REST
	s.downloadStatus.Store("downloadStatus", domain.Processing)
	go s.workerPoolDownloadTransactions(ctx, url)

	return domain.Processing, nil
}
