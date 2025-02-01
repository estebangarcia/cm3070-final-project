package api

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
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
	FileName    string
	DownloadURL string
	Digest      string
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

	w.Header().Set("Location", getBlobDownloadUrl(rh.Config.BaseURL, organization.Slug, registry.Slug, packageName, digest)+"?filename="+fileName)
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

		packages = append(packages, PythonPackage{
			FileName:    fileName,
			DownloadURL: downloadPythonPackageURL(rh.Config.BaseURL, organization.Slug, registry.Slug, packageName, fileName),
			Digest:      digest,
		})
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

	opts := oras.PackManifestOptions{
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
