package responses

type OrganizationStatsResponse struct {
	RegistryCount            int `json:"registry_count"`
	RepositoryCount          int `json:"repository_count"`
	ArtifactsCount           int `json:"artifacts_count"`
	StorageUsed              int `json:"storage_used"`
	VulnerableArtifactsCount int `json:"vulnerable_artifacts_count"`
}

type OrganizationMember struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}
