package api

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"maps"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/responses"
	"github.com/go-chi/chi/v5"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	oras_file "oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

const (
	PYTHON_ARTIFACT_TYPE                   string = "application/vnd.python.artifact"
	PYTHON_WHL_MEDIA_TYPE                  string = "application/vnd.python.whl.file"
	PYTHON_SIMPLE_REPOSITORY_TEMPLATE_NAME string = "simple.tmpl"
)

type PythonPackage struct {
	FileName         string
	DownloadURL      string
	Digest           string
	RequiresPython   string
	MetadataChecksum string
}

type PythonHandler struct {
	Config               *config.AppConfig
	ManifestRepository   *repositories.ManifestRepository
	RepositoryRepository *repositories.RepositoryRepository
}

// Handle the download of a python package by redirecting the client to the
// blob's URL location and specifing the filename that it should have when downloading
func (rh *PythonHandler) DownloadPackage(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	token := r.Context().Value("token").(string)
	packageName := chi.URLParam(r, "packageName")
	fileName := chi.URLParam(r, "fileName")
	digest := r.URL.Query().Get("digest")

	w.Header().Set("Location", getBlobDownloadUrl(rh.Config.GetBaseUrl(), organization.Slug, registry.Slug, packageName, digest)+"?filename="+fileName+"&token="+token)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// This build the simple repository index for python based on the uploaded manifests in the registry
func (rh *PythonHandler) SimpleRepositoryIndex(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	packageName := chi.URLParam(r, "packageName")

	lowerCasePackageName := strings.ToLower(packageName)
	if lowerCasePackageName != packageName {
		w.Header().Set("Location", simpleIndexUrl(rh.Config.GetBaseUrl(), organization.Slug, registry.Slug, lowerCasePackageName))
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	repo, found, err := rh.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, packageName)
	if !found {
		responses.OCIRepositoryUnknown(w, packageName, false)
		return
	}
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	// Get all OCI manifests that have the python artifact type
	manifests, err := rh.ManifestRepository.GetAllByTypeWithTags(r.Context(), PYTHON_ARTIFACT_TYPE, repo)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	var packages []PythonPackage

	// For each manifest build the data that will be used to render the index, including metadata that
	// is stored in the manifest's annotations
	for _, m := range manifests {
		var whlFileLayer *ent.ManifestLayer
		for _, l := range m.Edges.ManifestLayers {
			if l.MediaType == PYTHON_WHL_MEDIA_TYPE {
				whlFileLayer = l
			}
		}

		digest, _ := strings.CutPrefix(whlFileLayer.Digest, "sha256:")
		fileName := whlFileLayer.Annotations["org.opencontainers.image.title"]

		pkg := PythonPackage{
			FileName:    fileName,
			DownloadURL: downloadPythonPackageURL(rh.Config.GetBaseUrl(), organization.Slug, registry.Slug, packageName, fileName),
			Digest:      digest,
		}

		if requiresPython, ok := whlFileLayer.Annotations["Requires-Python"]; ok {
			pkg.RequiresPython = requiresPython
		}

		if metadataChecksum, ok := whlFileLayer.Annotations["Metadata-SHA256-Checksum"]; ok {
			pkg.MetadataChecksum = metadataChecksum
		}

		packages = append(packages, pkg)
	}

	html, err := rh.renderTemplate(packages)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

// Handle the upload of a python package, it parses the wheel's metadata and uses
// ORAS to upload the OCI artifact back to the registry
func (rh *PythonHandler) UploadPythonPackage(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	token := r.Context().Value("token").(string)
	pkgName := strings.ToLower(r.FormValue("name"))

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("content")
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	f, err := os.CreateTemp("/tmp", "python-pkg")
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	// Parse WHL metadata
	whlMetadata, err := parseWheelMetadata(f.Name())
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	store, err := oras_file.New("/tmp")
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	defer store.Close()

	// Create the layer descriptor for the python file

	layerDescriptor, err := store.Add(r.Context(), header.Filename, PYTHON_WHL_MEDIA_TYPE, f.Name())
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	// Store all WHL metadata in the layer's annotations
	maps.Copy(layerDescriptor.Annotations, whlMetadata)

	// Create the manifest, set the created timestamp to a fixed value
	// this is done to create a deterministic manifest digest for
	// proper manifest deduplication to avoid duplicate python packages to be uploadedd
	opts := oras.PackManifestOptions{
		ManifestAnnotations: map[string]string{
			"org.opencontainers.image.created": "2000-01-01T00:00:00Z",
		},
		Layers: []v1.Descriptor{
			layerDescriptor,
		},
	}
	// build the manifest descriptor
	manifestDescriptor, err := oras.PackManifest(r.Context(), store, oras.PackManifestVersion1_1, PYTHON_ARTIFACT_TYPE, opts)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	// Set the tag for the manifest, this is equal to the python package version
	pkgVersion := r.FormValue("version")
	if err = store.Tag(r.Context(), manifestDescriptor, pkgVersion); err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	// Upload the OCI artifact to the same server where this is being served from
	reg := "localhost:8081"
	repo, err := remote.NewRepository(reg + "/" + organization.Slug + "/" + registry.Slug + "/" + pkgName)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}
	repo.PlainHTTP = true

	// Set the credentials to be equal to the ones used by the client to initiate this request
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			AccessToken: token,
		}),
	}

	_, err = oras.Copy(r.Context(), store, pkgVersion, repo, pkgVersion, oras.DefaultCopyOptions)
	if err != nil {
		fmt.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Renders the simple index template
func (rh *PythonHandler) renderTemplate(packages []PythonPackage) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(PYTHON_SIMPLE_REPOSITORY_TEMPLATE_NAME).ParseFiles(wd + "/pkg/templates/python/" + PYTHON_SIMPLE_REPOSITORY_TEMPLATE_NAME)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, packages); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

// Build the download URL for a python package
func downloadPythonPackageURL(baseURL string, orgSlug string, registrySlug string, packageName string, fileName string) string {
	return fmt.Sprintf("%s/api/v1/%s/%s/python/simple/%s/%s", baseURL, orgSlug, registrySlug, packageName, fileName)
}

// Build the simple index url
func simpleIndexUrl(baseURL string, orgSlug string, registrySlug string, packageName string) string {
	return fmt.Sprintf("%v/api/v1/%v/%v/python/simple/%v", baseURL, orgSlug, registrySlug, strings.ToLower(packageName))
}

// Parse a wheel's metadata
func parseWheelMetadata(whlPath string) (map[string]string, error) {
	// Open the .whl (ZIP) file
	r, err := zip.OpenReader(whlPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Find the METADATA file inside the .whl archive
	var metadataFile *zip.File
	for _, file := range r.File {
		if strings.HasSuffix(file.Name, "METADATA") {
			metadataFile = file
			break
		}
	}

	if metadataFile == nil {
		return nil, fmt.Errorf("METADATA file not found in %s", whlPath)
	}

	// Read the METADATA file
	rc, err := metadataFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, rc)
	if err != nil {
		return nil, err
	}

	metadataChecksum := fmt.Sprintf("%x", sha256.Sum256(buffer.Bytes()))

	// A wheel metadata is the same format as an email's headers
	// use the mail library to parse it
	msg, err := mail.ReadMessage(&buffer)
	if err != nil {
		return nil, err
	}

	metadata := map[string]string{}
	for k, h := range msg.Header {
		if len(h) > 1 {
			jsonValue, err := json.Marshal(h)
			if err != nil {
				return nil, err
			}
			metadata[k] = string(jsonValue)
		} else {
			metadata[k] = h[0]
		}
	}

	msgBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}

	msgString := string(msgBytes)
	if len(msgString) > 0 {
		metadata["Description"] = msgString
	}

	metadata["Metadata-SHA256-Checksum"] = metadataChecksum

	return metadata, nil
}
