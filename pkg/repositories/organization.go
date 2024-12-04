package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
)

type OrganizationRepository struct {
	dbClient *ent.Client
}

func NewOrganizationRepository(dbClient *ent.Client) *OrganizationRepository {
	return &OrganizationRepository{
		dbClient: dbClient,
	}
}

func (orgRepo *OrganizationRepository) GetForUser(ctx context.Context, sub string) ([]*ent.Organization, error) {
	return orgRepo.dbClient.Organization.Query().Where(
		organization.HasMembersWith(
			user.Sub(sub),
		),
	).All(ctx)
}

func (orgRepo *OrganizationRepository) GetForUserAndSlug(ctx context.Context, sub string, orgSlug string) (*ent.Organization, bool, error) {
	org, err := orgRepo.dbClient.Organization.Query().Where(
		organization.And(
			organization.HasMembersWith(
				user.Sub(sub),
			),
			organization.Slug(orgSlug),
		),
	).First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return org, true, nil
}
