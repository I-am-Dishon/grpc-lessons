package main

import (
	"context"
	"fmt"
	calculatorpb "github.com/I-am-Dishon/grpc-lessons/grpc-calculator/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {}

func (*server) Add(ctx context.Context, req *calculatorpb.AddRequest) (*calculatorpb.AddResponse, error) {
	fmt.Printf("received Sum Rpc: %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	sum := firstNumber + secondNumber

	res := &calculatorpb.AddResponse{
		Result: sum,
	}
	return res, nil
}
func (*server) PrimeNumbers(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.AddService_PrimeNumbersServer) error {
	fmt.Printf("received Sum Rpc: %v\n", req)
	number := req.FirstNumber

	k := int64(2)
	for number > 1 {
		if number % k == 0 {
			//fmt.Printf("value : %v\n", k)
			stream.Send(&calculatorpb.PrimeNumberResponse{
				Result: k,
			})
			number = number/k
			fmt.Printf("k increased to: %v\n", k)
		} else {
			k++
		}

	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.AddService_ComputeAverageServer) error {
	fmt.Printf("received ComputeAverage Rpc \n")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}

		sum += req.Number
		count++

	}
	return nil
}

func main() {
	fmt.Println("Calculator Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterAddServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}