package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organizationinvite"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organizationmembership"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	ent_user "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/user"
)

type OrganizationInviteRepository struct {
}

func NewOrganizationInviteRepository() *OrganizationInviteRepository {
	return &OrganizationInviteRepository{}
}

func (orgRepo *OrganizationInviteRepository) HasInviteForOrganization(ctx context.Context, organization *ent.Organization, user *ent.User, email string) (bool, error) {
	dbClient := getClient(ctx)

	predicates := []predicate.OrganizationInvite{
		organizationinvite.OrganizationID(organization.ID),
	}

	if user != nil {
		predicates = append(predicates, organizationinvite.UserID(user.ID))
	} else {
		predicates = append(predicates, organizationinvite.Email(email))
	}

	return dbClient.OrganizationInvite.Query().Where(
		organizationinvite.And(
			predicates...,
		),
	).Exist(ctx)
}

func (orgRepo *OrganizationInviteRepository) InviteUserToOrganization(ctx context.Context, organization *ent.Organization, user *ent.User, email string, role string) error {
	dbClient := getClient(ctx)

	inviteCreate := dbClient.OrganizationInvite.Create().SetOrganization(organization).SetRole(organizationinvite.Role(role))
	if user != nil {
		inviteCreate.SetInvitee(user)
	} else {
		inviteCreate.SetEmail(email)
	}

	_, err := inviteCreate.Save(ctx)
	return err
}

func (orgRepo *OrganizationInviteRepository) GetInvitesForUser(ctx context.Context, userSub string) (ent.OrganizationInvites, error) {
	dbClient := getClient(ctx)

	return dbClient.OrganizationInvite.Query().Where(
		organizationinvite.HasInviteeWith(
			ent_user.SubEQ(userSub),
		),
	).WithOrganization().All(ctx)
}

func (orgRepo *OrganizationInviteRepository) HasInviteWithID(ctx context.Context, inviteId string, userSub string) (bool, error) {
	dbClient := getClient(ctx)

	return dbClient.OrganizationInvite.Query().Where(
		organizationinvite.InviteID(inviteId),
		organizationinvite.HasInviteeWith(
			ent_user.Sub(userSub),
		),
	).Exist(ctx)
}

func (orgRepo *OrganizationInviteRepository) FindInvitesForEmailAndLinkToUser(ctx context.Context, email string, user *ent.User) error {
	dbClient := getClient(ctx)

	_, err := dbClient.OrganizationInvite.Update().Where(
		organizationinvite.Email(email),
	).ClearEmail().SetUserID(user.ID).Save(ctx)
	return err
}

func (orgRepo *OrganizationInviteRepository) RejectInvite(ctx context.Context, inviteId string, userSub string) error {
	dbClient := getClient(ctx)

	_, err := dbClient.OrganizationInvite.Delete().Where(
		organizationinvite.InviteID(inviteId),
		organizationinvite.HasInviteeWith(
			ent_user.Sub(userSub),
		),
	).Exec(ctx)

	return err
}

func (orgRepo *OrganizationInviteRepository) AcceptInvite(ctx context.Context, inviteId string, userSub string) error {
	dbClient := getClient(ctx)

	invite, err := dbClient.OrganizationInvite.Query().Where(
		organizationinvite.InviteID(inviteId),
		organizationinvite.HasInviteeWith(
			ent_user.Sub(userSub),
		),
	).First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	_, err = dbClient.OrganizationMembership.Create().
		SetRole(organizationmembership.Role(invite.Role)).
		SetOrganizationID(invite.OrganizationID).
		SetUserID(*invite.UserID).
		Save(ctx)

	if err != nil {
		return err
	}

	return dbClient.OrganizationInvite.DeleteOneID(invite.ID).Exec(ctx)
}
