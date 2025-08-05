package main

import (
	"log"
	"net"
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ollamaClnt := embedding.NewOllamaClient("http://localhost:11434", "nomic-embed-text:v1.5")

	grpcServer := grpc.NewServer()
	s := &server{
		ollama: ollamaClnt,
	}
	pb.RegisterEmbedderServer(grpcServer, s)

	log.Println("Embedder gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
