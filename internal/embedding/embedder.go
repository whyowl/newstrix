package embedding

import (
	"context"
	"time"
)

type Embedder struct {
	client *EmbedClient
}

func NewEmbedder(addr string) (*Embedder, error) {
	cli, err := NewEmbedClient(addr)
	if err != nil {
		return nil, err
	}

	return &Embedder{client: cli}, nil
}

func (e *Embedder) Vectorize(ctx context.Context, text string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return e.client.Embed(ctx, text)
}
