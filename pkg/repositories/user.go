package repositories

import (
	"context"
	"fmt"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organizationmembership"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
	"github.com/gosimple/slug"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) GetUserBySub(ctx context.Context, sub string) (*ent.User, error) {
	dbClient := getClient(ctx)
	return dbClient.User.Query().Where(user.SubEQ(sub)).First(ctx)
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*ent.User, bool, error) {
	dbClient := getClient(ctx)
	user, err := dbClient.User.Query().Where(user.EmailEQ(email)).First(ctx)
	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, givenName string, familyName string, email string, sub string) (*ent.User, error) {
	dbClient := getClient(ctx)
	return dbClient.User.Create().SetGivenName(givenName).SetFamilyName(familyName).SetEmail(email).SetSub(sub).Save(ctx)
}

func (ur *UserRepository) CreateUserAndStartingOrg(ctx context.Context, givenName string, familyName string, email string, sub string) (*ent.User, *ent.Organization, error) {
	dbClient := getClient(ctx)

	user, err := dbClient.User.Create().SetGivenName(givenName).SetFamilyName(familyName).SetEmail(email).SetSub(sub).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	orgName := fmt.Sprintf("%s %s's Personal Organization", givenName, familyName)
	orgSlug := slug.Make(fmt.Sprintf("%s %s", givenName, familyName))

	orgSlugCount, err := dbClient.Organization.Query().Where(
		organization.SlugContains(orgSlug),
	).Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	if orgSlugCount > 0 {
		orgSlug = fmt.Sprintf("%v-%v", orgSlug, orgSlugCount)
	}

	org, err := dbClient.Organization.Create().SetName(orgName).SetSlug(orgSlug).SetIsPersonal(true).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	_, err = dbClient.OrganizationMembership.Create().SetOrganization(org).SetUser(user).SetRole(organizationmembership.RoleOwner).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	return user, org, nil
}
