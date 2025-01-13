package responses

type TagsListResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
