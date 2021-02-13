package main

import (
	"github.com/I-am-Dishon/grpc-lessons/grpc-calculator/calculatorpb"
	"context"
	"fmt"
	"net"
	"log"

	"google.golang.org/grpc"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *calculatorpb.Add) (*calculatorpb.AddResponse, error) {
	fmt.Printf("received Sum Rpc: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber

	res := &calculatorpb.AddResponse{
		SumResult: res
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterAddService(s, &server{})

}