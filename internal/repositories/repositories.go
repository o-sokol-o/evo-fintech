package repositories

import "github.com/jmoiron/sqlx"

type Repositories struct {
	RepoEVO    *RepoEVO
	RepoRemote *RepoRemote
}

func New(db *sqlx.DB) *Repositories {
	return &Repositories{
		RepoEVO:    NewRepoEVO(db),
		RepoRemote: NewRemoteRepo(db),
	}
}
