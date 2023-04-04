package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/dyjh/grpc_calculator/calculator"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Calculate(ctx context.Context, req *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	var res float32
	switch req.Operation {
	case pb.Operation_ADD:
		res = req.Num1 + req.Num2
	case pb.Operation_SUBTRACT:
		res = req.Num1 - req.Num2
	case pb.Operation_MULTIPLY:
		res = req.Num1 * req.Num2
	case pb.Operation_DIVIDE:
		res = req.Num1 / req.Num2
	default:
		return nil, fmt.Errorf("Invalid operation")
	}

	return &pb.CalculateResponse{Result: res}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, &server{})

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
