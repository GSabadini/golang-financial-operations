package adapter

import (
	"context"
	"encoding/json"
	"os"

	"go-opentelemetry-example/domain"
)

type (
	Transaction struct {
		client HTTPClient
		uri    string
	}

	transactionResponse struct {
		id int
	}
)

func NewTransaction(client HTTPClient, uri string) Transaction {
	return Transaction{
		client: client,
		uri:    uri,
	}
}

func (t Transaction) Transaction(_ context.Context, _ domain.TransactionInput) (int, error) {
	tr := &transactionResponse{}

	res, err := t.client.Get(os.Getenv("TRANSACTION_URI"))
	if err != nil {
		return tr.id, err
	}

	err = json.NewDecoder(res.Body).Decode(&tr)
	if err != nil {
		return tr.id, err
	}

	return tr.id, nil
}
