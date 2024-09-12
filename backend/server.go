package main

import (
	pb "backend/proto"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func (s *server) StreamTime(req *pb.TimeRequest, stream pb.MyService_StreamTimeServer) error {
	for {
		select {
		case <-stream.Context().Done():
			fmt.Println("StreamTime Done")
			return nil
		default:
			currentTime := time.Now().Format(time.RFC3339)
			if err := stream.Send(&pb.TimeResponse{CurrentTime: currentTime}); err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
			fmt.Println("StreamTime", currentTime)
		}
	}
}

func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMyServiceServer(grpcServer, &server{})

	// 서버 리플렉션 활성화
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	return grpcServer
}

func startHTTPServer(grpcServer *grpc.Server) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterMyServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		log.Printf("Received request with origin: %s", origin)
		return true // Allow all origins for testing
	}))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins for testing
		//AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"}, // Allow all headers for testing
		//AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposedHeaders:   []string{"Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("Received request: %s %s", req.Method, req.URL.Path)
		log.Printf("Request headers: %v", req.Header)

		if wrappedGrpc.IsGrpcWebRequest(req) {
			log.Println("Handling as gRPC-Web request")
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			log.Println("Handling as HTTP request")
			mux.ServeHTTP(resp, req)
		}

		log.Printf("Response headers: %v", resp.Header())
	}))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}

	log.Println("HTTP server listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	grpcServer := startGRPCServer()
	startHTTPServer(grpcServer)
}
