package api

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
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

func (rh *PythonHandler) DownloadPackage(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	packageName := chi.URLParam(r, "packageName")
	fileName := chi.URLParam(r, "fileName")
	digest := r.URL.Query().Get("digest")

	w.Header().Set("Location", getBlobDownloadUrl(rh.Config.GetBaseUrl(), organization.Slug, registry.Slug, packageName, digest)+"?filename="+fileName)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (rh *PythonHandler) SimpleRepositoryIndex(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	packageName := chi.URLParam(r, "packageName")

	repo, found, err := rh.RepositoryRepository.GetForRegistryByName(r.Context(), registry.ID, packageName)
	if !found {
		responses.OCIRepositoryUnknown(w, packageName, false)
		return
	}
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	manifests, err := rh.ManifestRepository.GetAllByTypeWithTags(r.Context(), PYTHON_ARTIFACT_TYPE, repo)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	var packages []PythonPackage

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
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

func (rh *PythonHandler) UploadPythonPackage(w http.ResponseWriter, r *http.Request) {
	organization := r.Context().Value("organization").(*ent.Organization)
	registry := r.Context().Value("registry").(*ent.Registry)
	token := r.Context().Value("token").(string)

	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("content")
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	f, err := os.CreateTemp("/tmp", "python-pkg")
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	whlMetadata, err := parseWheelMetadata(f.Name())
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	store, err := oras_file.New("/tmp")
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	defer store.Close()

	layerDescriptor, err := store.Add(r.Context(), header.Filename, PYTHON_WHL_MEDIA_TYPE, f.Name())
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	maps.Copy(layerDescriptor.Annotations, whlMetadata)

	opts := oras.PackManifestOptions{
		ManifestAnnotations: map[string]string{
			"org.opencontainers.image.created": "2000-01-01T00:00:00Z",
		},
		Layers: []v1.Descriptor{
			layerDescriptor,
		},
	}
	manifestDescriptor, err := oras.PackManifest(r.Context(), store, oras.PackManifestVersion1_1, PYTHON_ARTIFACT_TYPE, opts)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	pkgVersion := r.FormValue("version")
	if err = store.Tag(r.Context(), manifestDescriptor, pkgVersion); err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	pkgName := r.FormValue("name")

	reg := "localhost:8081"
	repo, err := remote.NewRepository(reg + "/" + organization.Slug + "/" + registry.Slug + "/" + pkgName)
	if err != nil {
		panic(err)
	}
	repo.PlainHTTP = true

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			AccessToken: token,
		}),
	}

	_, err = oras.Copy(r.Context(), store, pkgVersion, repo, pkgVersion, oras.DefaultCopyOptions)
	if err != nil {
		log.Println(err)
		responses.OCIInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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

func downloadPythonPackageURL(baseURL string, orgSlug string, registrySlug string, packageName string, fileName string) string {
	return fmt.Sprintf("%s/api/v1/%s/%s/python/simple/%s/%s", baseURL, orgSlug, registrySlug, packageName, fileName)
}

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
