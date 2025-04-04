package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	ent_organization "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
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

// Get all registries for the organization
func (registryRepo *RegistryRepository) GetForOrg(ctx context.Context, orgSlug string) ([]*ent.Registry, error) {
	dbClient := getClient(ctx)
	return dbClient.Registry.Query().Where(
		registry.HasOrganizationWith(
			organization.Slug(orgSlug),
		),
	).All(ctx)
}

// Create registry in the specified organization
func (registryRepo *RegistryRepository) CreateRegistry(ctx context.Context, name string, orgId int) (*ent.Registry, error) {
	registrySlug := slug.Make(name)
	dbClient := getClient(ctx)
	return dbClient.Registry.Create().SetOrganizationID(orgId).SetName(name).SetSlug(registrySlug).Save(ctx)
}

// Create registry by slug in the specified organization if user is an org member
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

// Get count of registries in the specified organization
func (registryRepo *RegistryRepository) GetCountForOrg(ctx context.Context, organization *ent.Organization) (int, error) {
	dbClient := getClient(ctx)

	return dbClient.Registry.Query().Where(
		registry.And(
			registry.HasOrganizationWith(
				ent_organization.ID(organization.ID),
			),
		),
	).Count(ctx)
}
