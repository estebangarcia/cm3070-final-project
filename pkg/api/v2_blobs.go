package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
)

type V2BlobsHandler struct {
	Config   *config.AppConfig
	S3Client *s3.Client
}

func (h *V2BlobsHandler) InitiateUploadSession(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName")
	log.Println(imageName)
	w.Header().Set("Location", fmt.Sprintf("%s/v2/hello/123", h.Config.BaseURL))
	w.WriteHeader(http.StatusAccepted)
}

func (h *V2BlobsHandler) HeadBlob(w http.ResponseWriter, r *http.Request) {
	imageName := r.Context().Value("imageName")
	blobDigest := r.Context().Value("digest")
	log.Println(imageName)
	log.Println(blobDigest)
	w.WriteHeader(http.StatusNotFound)
}
