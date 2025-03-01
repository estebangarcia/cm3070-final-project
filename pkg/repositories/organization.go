package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
)

type OrganizationRepository struct {
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{}
}

func (orgRepo *OrganizationRepository) GetForUser(ctx context.Context, sub string) ([]*ent.Organization, error) {
	dbClient := getClient(ctx)
	return dbClient.Organization.Query().Where(
		organization.HasMembersWith(
			user.Sub(sub),
		),
	).All(ctx)
}

func (orgRepo *OrganizationRepository) GetForUserAndSlug(ctx context.Context, sub string, orgSlug string) (*ent.Organization, bool, error) {
	dbClient := getClient(ctx)

	org, err := dbClient.Organization.Query().Where(
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
