package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
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
