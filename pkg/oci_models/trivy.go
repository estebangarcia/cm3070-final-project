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

type TrivyResult struct {
	Class           string
	Type            string
	Vulnerabilities []TrivyVulnerability
}

type TrivyReport struct {
	Results []TrivyResult
}
