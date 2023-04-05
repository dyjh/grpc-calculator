package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dyjh/grpc_calculator/calculator"
	"google.golang.org/grpc"
)

type server struct {
	calculator.UnimplementedCalculatorServer
}

func (*server) Calculate(ctx context.Context, req *calculator.CalculateRequest) (*calculator.CalculateResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	operation := req.GetOperation()

	var result float32
	switch operation {
	case calculator.Operation_ADD:
		result = num1 + num2
	case calculator.Operation_SUBTRACT:
		result = num1 - num2
	case calculator.Operation_MULTIPLY:
		result = num1 * num2
	case calculator.Operation_DIVIDE:
		if num2 == 0 {
			return nil, fmt.Errorf("division by zero is not allowed")
		}
		result = num1 / num2
	}

	return &calculator.CalculateResponse{Result: result}, nil
}

func (*server) Compare(ctx context.Context, req *calculator.CompareRequest) (*calculator.CompareResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()

	max := num1
	if num2 > num1 {
		max = num2
	}

	return &calculator.CompareResponse{Max: max}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
