package main

import (
	"context"
	"fmt"
	greetpb "github.com/I-am-Dishon/grpc-lessons/greet/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Greet Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaning(c)

	doBidiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Add unary rpc call")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Minions",
			SecondName: "Rain",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error when calling Greet grpc: %v", err)
	}

	log.Printf("Response from Greet service: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting {
			FirstName: "Minions",
			SecondName: "Rain",
		},
	}

	resStreams, err := c.GreetManyTimes(context.Background(), req)

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
		log.Printf("Response from GreetManyTimes service: %v", msg.GetResult())
	}

}

func doClientStreaning(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while sending stream %v", err)
	}

	// we iterate over  slice to send each message separately
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n ", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("unable to get LongGreet response from the server %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res.Result)
}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a biDi Streaming RPC")
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mary",
			},
		},
	}

	//create a stream by invoking client
	stream, err := c.GreetEveryone(context.Background())
	if err !=nil {
		log.Fatalf("error while creating stream: %v\n", err)
		return
	}

	waitc := make(chan  struct{})
	//send a bunch of messages to the client using go routine
	go func() {
		//functions to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending req: %v\n ", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		//functions to receive a bunch of messages
		for {

			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while sending stream %v", err)
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}