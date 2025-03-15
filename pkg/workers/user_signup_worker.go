package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
)

type CognitoUserSignUpEvent struct {
	Sub               string `json:"sub"`
	EmailVerified     string `json:"email_verified"`
	CognitoUserStatus string `json:"cognito:user_status"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

type UserSignupWorker struct {
	sqsClient                    *sqs.Client
	userRepository               *repositories.UserRepository
	organizationInviteRepository *repositories.OrganizationInviteRepository
}

func NewUserSignupWorker(sqsClient *sqs.Client, userRepository *repositories.UserRepository, organizationInviteRepository *repositories.OrganizationInviteRepository) *UserSignupWorker {
	return &UserSignupWorker{
		sqsClient:                    sqsClient,
		userRepository:               userRepository,
		organizationInviteRepository: organizationInviteRepository,
	}
}

func (w *UserSignupWorker) Handle(ctx context.Context, message types.Message) error {
	var userSignUpEvent CognitoUserSignUpEvent

	err := json.Unmarshal([]byte(*message.Body), &userSignUpEvent)
	if err != nil {
		fmt.Printf("error unmarshaling event %v\n", err)
		return err
	}

	user, _, err := w.userRepository.CreateUserAndStartingOrg(ctx, userSignUpEvent.GivenName, userSignUpEvent.FamilyName, userSignUpEvent.Email, userSignUpEvent.Sub)
	if err != nil {
		fmt.Printf("error creating user and org in database %v\n", err)
		return err
	}

	err = w.organizationInviteRepository.FindInvitesForEmailAndLinkToUser(ctx, user.Email, user)
	if err != nil {
		fmt.Printf("error linking invites to user %v\n", err)
		return err
	}

	return nil
}
