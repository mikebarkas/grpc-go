package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/mikebarkas/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	opts := grpc.WithInsecure()
	tls := true
	if tls {
		// ssl credentials
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Failed loading client certificates: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("created client: %f", c)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryDeadline(c, 6*time.Second)
	doUnaryDeadline(c, 2*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("start of doUnary\n")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "another",
			LastName:  "test",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling greet rpc %v", err)
	}
	log.Printf("response: %v\n", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "mike",
			LastName:  "nothing",
		},
	}

	stream, err := c.GreetMany(context.Background(), req)
	if err != nil {
		log.Fatalf("error calling greetmany")
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading stream")
		}
		log.Printf("Response from greetmany: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	// Example slice of request data.
	requests := []*greetpb.LongGreetRequest{

		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "one",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "two",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "three",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "four",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "five",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "six",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error calling LongGreet")
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response")
	}
	fmt.Printf("LongGreet response: %v ", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error creating stream")
	}

	// Example slice of request data.
	requests := []*greetpb.GreetEveryoneRequest{

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "one",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "two",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "three",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "four",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "five",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "six",
			},
		},
	}

	// Channel demonstration
	waitc := make(chan struct{})

	// Send messages with goroutine demonstration
	go func() {
		for _, req := range requests {
			fmt.Printf("sending request: %v \n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive messages in goroutine demonstration
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error receiving stream")
				break
			}
			fmt.Printf("Recieved: %v \n", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doUnaryDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Printf("start of doUnary with a deadline\n")

	req := &greetpb.GreetDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "another",
			LastName:  "test",
		},
	}

	// Create a context and include a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit")
			} else {
				fmt.Printf("an unexpected error: %v \n", statusErr.Message())
			}
		} else {
			log.Fatalf("error calling greet with deadline rpc: %v", err)
		}
		return
	}
	log.Printf("response: %v\n", res.Result)
}
