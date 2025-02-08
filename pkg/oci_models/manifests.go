package oci_models

type OCIV1Annotations map[string]string

type OCIV1ManifestBlobRef struct {
	MediaType           string           `json:"mediaType"`
	Digest              string           `json:"digest"`
	Size                int              `json:"size"`
	Data                *string          `json:"data,omitempty"`
	Annotations         OCIV1Annotations `json:"annotations,omitempty"`
	NewUnspecifiedField *string          `json:"newUnspecifiedField,omitempty"`
}

type OCIV1Manifest struct {
	SchemaVersion int                    `json:"schemaVersion"`
	MediaType     string                 `json:"mediaType"`
	Digest        *string                `json:"digest,omitempty"`
	Config        OCIV1ManifestBlobRef   `json:"config"`
	Layers        []OCIV1ManifestBlobRef `json:"layers"`
	Subject       *OCIV1ManifestBlobRef  `json:"subject,omitempty"`
	Annotations   OCIV1Annotations       `json:"annotations,omitempty"`
	ArtifactType  *string                `json:"artifactType,omitempty"`
}

type OCIV1ManifestIndex struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Manifests     []OCIV1Manifest `json:"manifests"`
}

func NewOCIV1ManifestIndex(manifests []OCIV1Manifest) OCIV1ManifestIndex {
	return OCIV1ManifestIndex{
		SchemaVersion: 2,
		MediaType:     "application/vnd.oci.image.index.v1+json",
		Manifests:     manifests,
	}
}
