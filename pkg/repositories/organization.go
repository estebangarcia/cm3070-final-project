package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organizationmembership"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
	"github.com/gosimple/slug"
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

func (orgRepo *OrganizationRepository) CreateOrganizationWithOwner(ctx context.Context, user *ent.User, orgName string) (*ent.Organization, error) {
	dbClient := getClient(ctx)

	orgSlug := slug.Make(orgName)

	org, err := dbClient.Organization.Create().SetName(orgName).SetSlug(orgSlug).SetIsPersonal(false).Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = dbClient.OrganizationMembership.Create().SetOrganization(org).SetUser(user).SetRole(organizationmembership.RoleOwner).Save(ctx)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (orgRepo *OrganizationRepository) GetOrganizationMembers(ctx context.Context, organization *ent.Organization) (ent.Users, error) {
	dbClient := getClient(ctx)
	return dbClient.User.Query().Where(
		user.HasJoinedOrganizationsWith(
			organizationmembership.OrganizationID(organization.ID),
		),
	).WithJoinedOrganizations(
		func(omq *ent.OrganizationMembershipQuery) {
			omq.Where(
				organizationmembership.OrganizationID(organization.ID),
			)
		},
	).All(ctx)
}
