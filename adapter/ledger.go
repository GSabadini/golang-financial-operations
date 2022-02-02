package adapter

import (
	"context"
	"encoding/json"
	"go-opentelemetry-example/domain"
	"os"
)

type (
	Ledger struct {
		client HTTPClient
		uri    string
	}

	ledgerResponse struct {
		id int
	}
)

func NewLedger(client HTTPClient, uri string) Ledger {
	return Ledger{
		client: client,
		uri:    uri,
	}
}

func (l Ledger) AvailableEntries(context.Context, domain.LedgerInput) error {
	lr := &ledgerResponse{}

	res, err := l.client.Get(os.Getenv("LEDGER_URI"))
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&lr)
	if err != nil {
		return err
	}

	return nil
}
