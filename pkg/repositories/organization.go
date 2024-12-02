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
