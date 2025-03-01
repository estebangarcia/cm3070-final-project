package oci_models

type TrivyCVSS struct {
	V2Vector string  `json:"V2Vector"`
	V3Vector string  `json:"V3Vector"`
	V2Score  float32 `json:"V2Score"`
	V3Score  float32 `json:"V3Score"`
}

type TrivyVulnerability struct {
	VulnerabilityID  string               `json:"VulnerabilityID"`
	PackageName      string               `json:"PkgName"`
	PackageID        string               `json:"PkgID"`
	InstalledVersion string               `json:"InstalledVersion"`
	FixedVersion     string               `json:"FixedVersion"`
	Status           string               `json:"Status"`
	PrimaryURL       string               `json:"PrimaryURL"`
	Title            string               `json:"Title"`
	Severity         string               `json:"Severity"`
	CVSS             map[string]TrivyCVSS `json:"CVSS"`
}

type TrivyMisconfiguration struct {
	ID         string `json:"ID"`
	Title      string `json:"Title"`
	Message    string `json:"Message"`
	Resolution string `json:"Resolution"`
	PrimaryURL string `json:"PrimaryURL"`
	Severity   string `json:"Severity"`
}

type TrivyResult struct {
	Target            string
	Class             string
	Type              string
	Vulnerabilities   []TrivyVulnerability
	Misconfigurations []TrivyMisconfiguration
}

type TrivyReport struct {
	Results []TrivyResult
}
