package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/requests"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/go-chi/chi/v5"
)

type OrganizationsHandler struct {
	Config                       *config.AppConfig
	OrganizationRepository       *repositories.OrganizationRepository
	OrganizationInviteRepository *repositories.OrganizationInviteRepository
	UserRepository               *repositories.UserRepository
	RegistryRepository           *repositories.RegistryRepository
	RepositoryRepository         *repositories.RepositoryRepository
	ManifestRepository           *repositories.ManifestRepository
	SESClient                    *sesv2.Client
}

// Get all the organizations for the user
func (oh *OrganizationsHandler) GetOrganizationsForUser(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)

	orgs, err := oh.OrganizationRepository.GetForUser(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orgs)
}

// Get a specific organization by its slug
func (oh *OrganizationsHandler) GetOrganizationsBySlugForUser(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)
	orgSlug := chi.URLParam(r, "organizationSlug")

	org, found, err := oh.OrganizationRepository.GetForUserAndSlug(r.Context(), userSub, orgSlug)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if !found {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

// Create an organization
func (oh *OrganizationsHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	userSub := r.Context().Value("user_sub").(string)

	createOrgRequest, err := requests.BindRequest[requests.CreateOrganizationRequest](r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	user, err := oh.UserRepository.GetUserBySub(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	org, err := oh.OrganizationRepository.CreateOrganizationWithAdmin(r.Context(), user, createOrgRequest.Name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(org)
}

// Get the organization stats, to be shown in the UI dashboard
func (oh *OrganizationsHandler) GetOrganizationStats(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)

	storageUsed, err := oh.ManifestRepository.GetStorageUsedInBytesForOrganization(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	registryCount, err := oh.RegistryRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	repositoryCount, err := oh.RepositoryRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	artifactCount, err := oh.ManifestRepository.GetCountForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	artifactWithVulnerabilitiesCount, err := oh.ManifestRepository.GetCountWithVulnerabilitiesForOrg(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	response := &responses.OrganizationStatsResponse{
		RegistryCount:            registryCount,
		RepositoryCount:          repositoryCount,
		ArtifactsCount:           artifactCount,
		StorageUsed:              storageUsed,
		VulnerableArtifactsCount: artifactWithVulnerabilitiesCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get all the members of an organization
func (oh *OrganizationsHandler) GetOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)

	members, err := oh.OrganizationRepository.GetOrganizationMembers(r.Context(), organization)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	var membersResponse []responses.OrganizationMember

	for _, member := range members {
		membersResponse = append(membersResponse, responses.OrganizationMember{
			GivenName:  member.GivenName,
			FamilyName: member.FamilyName,
			Email:      member.Email,
			Role:       member.Edges.JoinedOrganizations[0].Role.String(),
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(membersResponse)
}

// Invite a new user to the organization
func (oh *OrganizationsHandler) InviteToOrganization(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	userSub := r.Context().Value("user_sub").(string)

	user, err := oh.UserRepository.GetUserBySub(r.Context(), userSub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	inviteToOrgRequest, err := requests.BindRequest[requests.InviteToOrganization](r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(400)
		return
	}

	// Get the user by email in the invite, the user might not exist
	inviteeUser, inviteeFound, err := oh.UserRepository.GetUserByEmail(r.Context(), inviteToOrgRequest.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	// Check if the user already has an invite for this organization
	hasInvite, err := oh.OrganizationInviteRepository.HasInviteForOrganization(r.Context(), organization, inviteeUser, inviteToOrgRequest.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	if hasInvite {
		w.WriteHeader(400)
		return
	}

	// Create the invitation for the user. If the user doesn't exist this will add the user's email to the database instead of
	// creating a relationship to an existing user.
	if err := oh.OrganizationInviteRepository.InviteUserToOrganization(r.Context(), organization, inviteeUser, inviteToOrgRequest.Email, inviteToOrgRequest.Role); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	// Send email to notify user of their inviation
	if err := oh.sendInviteEmail(context.Background(), user.GivenName+" "+user.FamilyName, inviteToOrgRequest.Email, organization.Name, inviteeFound); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Build the email's body and use AWS SES for sending the email
func (oh *OrganizationsHandler) sendInviteEmail(ctx context.Context, inviterName string, email string, orgName string, userExists bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpl, err := template.New("organization_invite.tmpl").Funcs(template.FuncMap{
		"arr": func(els ...any) []any { return els },
	}).ParseGlob(wd + "/pkg/templates/emails/*.tmpl")
	if err != nil {
		return err
	}

	type EmailData struct {
		OrganizationName string
		InviterName      string
		AcceptLink       string
		SignupLink       string
		UserExists       bool
	}

	emailData := EmailData{
		OrganizationName: orgName,
		InviterName:      inviterName,
		UserExists:       userExists,
		AcceptLink:       oh.Config.FrontendBaseURL + "/invites/accept",
		SignupLink:       oh.Config.FrontendBaseURL,
	}

	var tpl bytes.Buffer

	if err := tmpl.Execute(&tpl, emailData); err != nil {
		return err
	}

	_, err = oh.SESClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(oh.Config.SES.FromEmailAddress),
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(fmt.Sprintf("You've been invited to the %v organization", orgName)),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(tpl.String()),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
