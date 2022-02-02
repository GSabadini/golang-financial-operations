package adapter

import (
	"context"
	"encoding/json"
	"go-opentelemetry-example/domain"
	"os"
)

type (
	RulesEvaluator struct {
		client HTTPClient
		uri    string
	}

	rulesEvaluatorResponse struct {
		id int
	}
)

func NewRulesEvaluator(client HTTPClient, uri string) RulesEvaluator {
	return RulesEvaluator{
		client: client,
		uri:    uri,
	}
}

func (r RulesEvaluator) SpendingControls(_ context.Context, _ domain.RulesInput) error {
	rer := &rulesEvaluatorResponse{}

	res, err := r.client.Get(os.Getenv("LEDGER_URI"))
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&rer)
	if err != nil {
		return err
	}

	return nil
}
