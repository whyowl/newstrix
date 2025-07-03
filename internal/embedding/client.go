package embedding

import (
	"context"
	"google.golang.org/grpc"
	"newstrix/internal/embedding/proto" // путь до автогенерированного кода
)

type EmbedClient struct {
	client pb.EmbedderClient
}

func NewEmbedClient(addr string) (*EmbedClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &EmbedClient{
		client: pb.NewEmbedderClient(conn),
	}, nil
}

func (ec *EmbedClient) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := ec.client.Embed(ctx, &pb.EmbedRequest{Text: text})
	if err != nil {
		return nil, err
	}
	return resp.Vector, nil
}
