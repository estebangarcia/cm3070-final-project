package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
	"github.com/gosimple/slug"
)

type RegistryRepository struct {
	dbClient *ent.Client
}

func NewRegistryRepository(dbClient *ent.Client) *RegistryRepository {
	return &RegistryRepository{
		dbClient: dbClient,
	}
}

func (registryRepo *RegistryRepository) GetForOrg(ctx context.Context, orgSlug string) ([]*ent.Registry, error) {
	return registryRepo.dbClient.Registry.Query().Where(
		registry.HasOrganizationWith(
			organization.Slug(orgSlug),
		),
	).All(ctx)
}

func (registryRepo *RegistryRepository) CreateRegistry(ctx context.Context, name string, orgId int) (*ent.Registry, error) {
	registrySlug := slug.Make(name)
	return registryRepo.dbClient.Registry.Create().SetOrganizationID(orgId).SetName(name).SetSlug(registrySlug).Save(ctx)
}

func (registryRepo *RegistryRepository) GetForOrgAndUser(ctx context.Context, userSub string, orgSlug string, registrySlug string) (*ent.Registry, bool, error) {
	registry, err := registryRepo.dbClient.Registry.Query().Where(
		registry.And(
			registry.Slug(registrySlug),
			registry.HasOrganizationWith(
				organization.And(
					organization.Slug(orgSlug),
					organization.HasMembersWith(
						user.Sub(userSub),
					),
				),
			),
		),
	).WithOrganization().First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return registry, true, nil
}
