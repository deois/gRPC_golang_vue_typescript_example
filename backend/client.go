package main

import (
	pb "backend/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	client := pb.NewMyServiceClient(conn)
	ctx := context.Background()
	req := &pb.HelloRequest{Name: "World"}
	resp, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Response: %s", resp.Message)

	req.Name = "Gopher"
	resp, err = client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Response: %s", resp.Message)
	fmt.Println("Done")
}
