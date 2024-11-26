package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

type RepositoryRepository struct {
	dbClient *ent.Client
}

func NewRepositoryRepository(dbClient *ent.Client) *RepositoryRepository {
	return &RepositoryRepository{
		dbClient: dbClient,
	}
}

func (rr *RepositoryRepository) GetOrCreateRepository(ctx context.Context, repositoryName string) (*ent.Repository, error) {
	repository, err := rr.dbClient.Repository.Query().Where(repository.Name(repositoryName)).First(ctx)
	if err != nil && ent.IsNotFound(err) {
		repository, err = rr.dbClient.Repository.Create().SetName(repositoryName).Save(ctx)
	}
	if err != nil {
		return nil, err
	}

	return repository, err
}
