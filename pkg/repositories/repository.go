package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
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

func (rr *RepositoryRepository) GetForRegistryByName(ctx context.Context, registryId int, repositoryName string) (*ent.Repository, bool, error) {
	repo, err := rr.dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.ID(registryId),
			),
			repository.Name(repositoryName),
		),
	).First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return repo, true, nil
}

func (rr *RepositoryRepository) GetAllForRegistry(ctx context.Context, registryId int) ([]*ent.Repository, error) {
	return rr.dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.ID(registryId),
			),
		),
	).All(ctx)
}

func (rr *RepositoryRepository) GetOrCreateRepository(ctx context.Context, registryId int, repositoryName string) (*ent.Repository, error) {
	repository, err := rr.dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.ID(registryId),
			),
			repository.Name(repositoryName),
		),
	).First(ctx)
	if err != nil && ent.IsNotFound(err) {
		repository, err = rr.dbClient.Repository.Create().SetName(repositoryName).SetRegistryID(registryId).Save(ctx)
	}
	if err != nil {
		return nil, err
	}

	return repository, err
}
