package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	ent_organization "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

type RepositoryRepository struct {
}

func NewRepositoryRepository() *RepositoryRepository {
	return &RepositoryRepository{}
}

func (rr *RepositoryRepository) GetForRegistryByName(ctx context.Context, registryId int, repositoryName string) (*ent.Repository, bool, error) {
	dbClient := getClient(ctx)
	repo, err := dbClient.Repository.Query().Where(
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
	dbClient := getClient(ctx)
	return dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.ID(registryId),
			),
		),
	).All(ctx)
}

func (rr *RepositoryRepository) GetOrCreateRepository(ctx context.Context, registryId int, repositoryName string) (*ent.Repository, error) {
	dbClient := getClient(ctx)
	repository, err := dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.ID(registryId),
			),
			repository.Name(repositoryName),
		),
	).First(ctx)
	if err != nil && ent.IsNotFound(err) {
		repository, err = dbClient.Repository.Create().SetName(repositoryName).SetRegistryID(registryId).Save(ctx)
	}
	if err != nil {
		return nil, err
	}

	return repository, err
}

func (rr *RepositoryRepository) GetCountForOrg(ctx context.Context, organization *ent.Organization) (int, error) {
	dbClient := getClient(ctx)
	return dbClient.Repository.Query().Where(
		repository.And(
			repository.HasRegistryWith(
				registry.HasOrganizationWith(
					ent_organization.ID(organization.ID),
				),
			),
		),
	).Count(ctx)
}
