package main

import (
	pb "backend/proto"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
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
			fmt.Println("StreamTime")
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
	e := echo.New()

	// 미들웨어 설정
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"*"},
	}))

	// gRPC 핸들러 설정
	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		log.Printf("Received request with origin: %s", origin)
		return true // Allow all origins for testing
	}))

	// Define and initialize mux
	mux := runtime.NewServeMux()
	err := pb.RegisterMyServiceHandlerServer(context.Background(), mux, &server{})
	if err != nil {
		log.Fatalf("Failed to register gRPC handler: %v", err)
	}

	e.Any("/*", func(c echo.Context) error {
		req := c.Request()
		resp := c.Response()

		if wrappedGrpc.IsGrpcWebRequest(req) {
			log.Println("Handling as gRPC-Web request")
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			log.Println("Handling as HTTP request")
			mux.ServeHTTP(resp, req)
		}

		return nil
	})

	// HTTP 서버 시작
	log.Println("HTTP server listening on :8080")
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	grpcServer := startGRPCServer()
	startHTTPServer(grpcServer)
}
