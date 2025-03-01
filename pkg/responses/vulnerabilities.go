package responses

import "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"

type VulnerabilitiesResponse struct {
	Vulnerabilities   ent.Vulnerabilities           `json:"vulnerabilities"`
	Misconfigurations ent.ManifestMisconfigurations `json:"misconfigurations"`
}
