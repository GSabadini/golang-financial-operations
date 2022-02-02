package adapter

import (
	"context"
	"go-opentelemetry-example/domain"
)

type Stream struct{}

func NewStream() Stream {
	return Stream{}
}

func (s Stream) Publish(_ context.Context, _ domain.StreamInput) {

}
