package main

import (
	"log"
	"net"
	"newstrix/internal/config"
	"newstrix/internal/embedding"

	"context"
	"google.golang.org/grpc"
	"newstrix/internal/embedding/proto"
)

type server struct {
	pb.UnimplementedEmbedderServer
	ollama *embedding.OllamaClient
}

func (s *server) Embed(ctx context.Context, req *pb.EmbedRequest) (*pb.EmbedResponse, error) {
	vector, err := s.ollama.Embed(ctx, req.Text)
	if err != nil {
		return nil, err // TODO hide error to out format
	}
	return &pb.EmbedResponse{Vector: vector}, nil
}

func main() {

	cfg := config.Load()

	lis, err := net.Listen("tcp", cfg.GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ollamaClnt := embedding.NewOllamaClient(cfg.OllamaURL, cfg.OllamaModel)

	grpcServer := grpc.NewServer()
	s := &server{
		ollama: ollamaClnt,
	}
	pb.RegisterEmbedderServer(grpcServer, s)

	log.Printf("Embedder gRPC server running on %s \n", cfg.GrpcAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
