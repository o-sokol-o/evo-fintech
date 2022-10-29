package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

const tableSourceCSV = "sourceCSV"

type RepoRemote struct {
	db *sqlx.DB
}

func NewRemoteRepo(db *sqlx.DB) *RepoRemote {
	return &RepoRemote{db: db}
}

func (r *RepoRemote) Get(ctx context.Context) ([]domain.Transaction, error) {
	var sourceCSV []domain.Transaction

	if err := r.db.SelectContext(ctx, &sourceCSV, `select * from `+tableSourceCSV); err != nil {
		return nil, err
	}

	return sourceCSV, nil
}
