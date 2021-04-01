package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/mikebarkas/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Unary request and response.
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func was invoked with: %v\n", req)

	res := &greetpb.GreetResponse{
		Result: &greetpb.Greeting{
			FirstName: req.Greeting.GetFirstName(),
			LastName:  req.Greeting.GetLastName(),
		},
	}
	return res, nil
}

// Server stream response.
func (*server) GreetMany(req *greetpb.GreetManyRequest, stream greetpb.GreetService_GreetManyServer) error {
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// Client stream request.
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Finished receiving, now sending response \n")
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error reading client stream %v", err)
		}
		//firstName := req.Greeting.GetFirstName()
		result += " Hello " + req.Greeting.GetFirstName() + "\n "
	}

}

// Bidirectional request response stream.
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Finished sending GreetEveryone")
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading client stream %v", err)
		}

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + req.Greeting.GetFirstName() + " \n",
		})
		if sendErr != nil {
			log.Fatalf("Error sending stream %v", err)
		}
	}
}

func main() {
	fmt.Println("Server main")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failted to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
