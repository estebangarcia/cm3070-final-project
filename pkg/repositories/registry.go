package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
	"github.com/gosimple/slug"
)

type RegistryRepository struct {
}

func NewRegistryRepository() *RegistryRepository {
	return &RegistryRepository{}
}

func (registryRepo *RegistryRepository) GetForOrg(ctx context.Context, orgSlug string) ([]*ent.Registry, error) {
	dbClient := getClient(ctx)
	return dbClient.Registry.Query().Where(
		registry.HasOrganizationWith(
			organization.Slug(orgSlug),
		),
	).All(ctx)
}

func (registryRepo *RegistryRepository) CreateRegistry(ctx context.Context, name string, orgId int) (*ent.Registry, error) {
	registrySlug := slug.Make(name)
	dbClient := getClient(ctx)
	return dbClient.Registry.Create().SetOrganizationID(orgId).SetName(name).SetSlug(registrySlug).Save(ctx)
}

func (registryRepo *RegistryRepository) GetForOrgAndUser(ctx context.Context, userSub string, orgSlug string, registrySlug string) (*ent.Registry, bool, error) {
	dbClient := getClient(ctx)

	orgPredicate := []predicate.Organization{
		organization.Slug(orgSlug),
	}

	if userSub != "" {
		orgPredicate = append(orgPredicate, organization.HasMembersWith(
			user.Sub(userSub),
		))
	}

	registry, err := dbClient.Registry.Query().Where(
		registry.And(
			registry.Slug(registrySlug),
			registry.HasOrganizationWith(
				organization.And(
					orgPredicate...,
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
