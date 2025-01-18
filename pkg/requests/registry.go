package requests

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CreateRegistryRequest struct {
	Name string `json:"name" validate:"required"`
}

func BindRequest[T any](r *http.Request) (*T, error) {
	var requestData T

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		return nil, err
	}

	err = validate.Struct(requestData)
	return &requestData, err
}

func BindRequestFromBytes[T any](data []byte) (*T, error) {
	var requestData T

	validate := validator.New(validator.WithRequiredStructEnabled())

	reader := bytes.NewReader(data)
	err := json.NewDecoder(reader).Decode(&requestData)
	if err != nil {
		return nil, err
	}

	err = validate.Struct(requestData)
	return &requestData, err
}
