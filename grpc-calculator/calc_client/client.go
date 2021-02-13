package main

import (
	"context"
	"fmt"
	calculatorpb "github.com/I-am-Dishon/grpc-lessons/grpc-calculator/proto"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewAddServiceClient(cc)

	doUnary(c)

	doPrimeNumber(c)

	doComputeAverage(c)
}

func doUnary(c calculatorpb.AddServiceClient) {
	fmt.Println("Starting Add unary rpc call")

	req := &calculatorpb.AddRequest{
		FirstNumber: 34,
		SecondNumber: 56,
	}

	res, err := c.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("error when calling Add grpc: %v", err)
	}

	log.Printf("Response from add service: %v", res.Result)
}

func doPrimeNumber(c calculatorpb.AddServiceClient) {
	fmt.Println("Starting to do a PrimeNumber Server Streaming RPC")

	req := &calculatorpb.PrimeNumberRequest{
		FirstNumber: 120,
	}

	resStreams, err := c.PrimeNumbers(context.Background(), req)

	if err != nil {
		log.Fatalf("error when calling  GreetManyTimes service: %v", err)
	}

	for {
		msg, err := resStreams.Recv()
		if err == io.EOF {
			// we've reached end the end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from PrimeNumber service: %v", msg.GetResult())
	}
}

func doComputeAverage(c calculatorpb.AddServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while openning stream %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number :=  range numbers{
		fmt.Printf("sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v\n", err)
	}

	fmt.Printf("The average is: %v\n", res.Average)
}