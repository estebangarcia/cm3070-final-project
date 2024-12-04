package repositories

import (
	"context"
	"fmt"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/gosimple/slug"
)

type UserRepository struct {
	dbClient *ent.Client
}

func NewUserRepository(dbClient *ent.Client) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, givenName string, familyName string, email string, sub string) (*ent.User, error) {
	return ur.dbClient.User.Create().SetGivenName(givenName).SetFamilyName(familyName).SetEmail(email).SetSub(sub).Save(ctx)
}

func (ur *UserRepository) CreateUserAndStartingOrg(ctx context.Context, givenName string, familyName string, email string, sub string) (*ent.User, *ent.Organization, error) {
	tx, err := ur.dbClient.Tx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	user, err := tx.User.Create().SetGivenName(givenName).SetFamilyName(familyName).SetEmail(email).SetSub(sub).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	orgName := fmt.Sprintf("%s %s's Personal Organization", givenName, familyName)
	orgSlug := slug.Make(fmt.Sprintf("%s %s", givenName, familyName))

	org, err := tx.Organization.Create().SetName(orgName).SetSlug(orgSlug).SetIsPersonal(true).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.OrganizationMembership.Create().SetOrganization(org).SetUser(user).SetRole(0).Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}

	return user, org, nil
}
