package main

import (
	"context"
	"log"
	"net"

	pb "eduData/grpc"

	"google.golang.org/grpc"
)

// server 是 AuthService 的实现
type server struct {
	pb.UnimplementedAuthServiceServer
}

// Signin 实现
func (s *server) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	// 在这里实现您的登录逻辑
	// 这是一个示例响应
	return &pb.SigninResponse{
		CookieJar:    []byte("example_cookie"),
		ErrorMessage: "",
		Success:      true,
	}, nil
}

// GetData 实现
func (s *server) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	// 在这里实现您的获取数据逻辑
	// 这是一个示例响应
	return &pb.GetDataResponse{
		Data:         []byte("example_data"),
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
