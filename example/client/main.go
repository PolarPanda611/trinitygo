package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	helloworldpb "trinitygo/example/pb/helloworld"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:50051"
	defaultName = "test trinity go "
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworldpb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for i := 0; i < 100; i++ {
		go func() {
			md := metadata.Pairs("trace_id", uuid.New().String())
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			r, err := c.SayHello(ctx, &helloworldpb.HelloRequest{Name: name})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(r.GetMessage())
			}

		}()

	}
	time.Sleep(time.Second * 1)

}
