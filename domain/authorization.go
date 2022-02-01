package domain

import "errors"

var (
	ErrDuplicatedAuthorization = errors.New("duplicated authorization")
)

type Authorization struct {
	id             string
	idempotenceKey string
}

func (a Authorization) AlreadyExists() bool {
	return a.id != ""
}

func (a Authorization) ID() string {
	return a.id
}

func (a Authorization) IdempotenceKey() string {
	return a.idempotenceKey
}
