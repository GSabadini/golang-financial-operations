package domain

type (
	LedgerInput struct {
		operation      string
		currency       string
		amount         int
		accountID      int
		idempotenceKey string
	}

	RulesInput struct {
		operation      string
		currency       string
		amount         int
		accountID      int
		idempotenceKey string
	}

	AuthorizationInput struct {
		operation      string
		currency       string
		amount         int
		accountID      int
		idempotenceKey string
	}

	TransactionInput struct {
		authorizationID int
		amount          int
		currency        string
		operation       string
		descriptor      string
		accountID       int
	}

	StreamInput struct {
		authorizationID int
		transactionID   int
		amount          int
		currency        string
		operation       string
		descriptor      string
		accountID       int
	}
)
