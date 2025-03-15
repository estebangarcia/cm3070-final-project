package requests

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required"`
}

type InviteToOrganization struct {
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required,oneof=owner manager member"`
}
