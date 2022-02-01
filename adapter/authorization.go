package adapter

import (
	"context"
	"encoding/json"
	"go-opentelemetry-example/domain"
	"os"
)

type (
	Authorization struct {
		client HTTPClient
	}

	authorizationResponse struct {
		id string
	}
)

func (a Authorization) Create(_ context.Context, _ string, _ int) error {
	res, err := a.client.Get(os.Getenv("AUTHORIZATION_URI"))
	if err != nil {
		return err
	}

	b := &authorizationResponse{}
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return err
	}

	return nil
}

func (a Authorization) FindByIdempotenceKey(_ context.Context, _ string) (domain.Authorization, error) {
	res, err := a.client.Get(os.Getenv("AUTHORIZATION_URI"))
	if err != nil {
		return domain.Authorization{}, err
	}

	b := &authorizationResponse{}
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return domain.Authorization{}, err
	}

	return domain.Authorization{}, nil
}
