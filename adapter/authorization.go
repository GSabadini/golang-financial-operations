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
		uri    string
	}

	authorizationResponse struct {
		id int
	}
)

func NewAuthorization(client HTTPClient, uri string) Authorization {
	return Authorization{
		client: client,
		uri:    uri,
	}
}

func (a Authorization) Create(_ context.Context, _ domain.AuthorizationInput) (int, error) {
	b := &authorizationResponse{}

	res, err := a.client.Get(a.uri)
	if err != nil {
		return b.id, err
	}

	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return b.id, err
	}

	return b.id, nil
}

func (a Authorization) FindByIdempotenceKey(_ context.Context, _ domain.AuthorizationInput) (domain.Authorization, error) {
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
