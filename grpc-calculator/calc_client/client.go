package main

import (
	"fmt"
	"github.com/I-am-Dishon/grpc-lessons/grpc-calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewAddService(cc)

	doUnary(c)
}

func doUnary(c calculatorpb.AddService) {
	fmt.Println("Starting Add unary rpc call")

	req := &calculatorpb.AddService{
		AddRequest: &AddRequest{
			FirstNumber: 1,
			SecondNumber: 2,
		}
	}

	res, err := c.Add(context.BackGround(), req)
	if err != nil {
		log.Fatalf("error when calling Add grpc: %v", err)
	}

	log.Printf("Response from ADD service: %v", res.Result)
}