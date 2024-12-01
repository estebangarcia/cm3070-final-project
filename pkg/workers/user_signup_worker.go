package workers

import (
	"context"
	"encoding/json"
	"log"

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
	sqsClient      *sqs.Client
	userRepository *repositories.UserRepository
}

func NewUserSignupWorker(sqsClient *sqs.Client, userRepository *repositories.UserRepository) *UserSignupWorker {
	return &UserSignupWorker{
		sqsClient:      sqsClient,
		userRepository: userRepository,
	}
}

func (w *UserSignupWorker) Handle(ctx context.Context, message types.Message) error {
	var userSignUpEvent CognitoUserSignUpEvent

	err := json.Unmarshal([]byte(*message.Body), &userSignUpEvent)
	if err != nil {
		log.Printf("error unmarshaling event %v", err)
		return err
	}

	_, err = w.userRepository.CreateUser(ctx, userSignUpEvent.GivenName, userSignUpEvent.FamilyName, userSignUpEvent.Email, userSignUpEvent.Sub)
	if err != nil {
		log.Printf("error creating user in database %v", err)
		return err
	}

	return nil
}
