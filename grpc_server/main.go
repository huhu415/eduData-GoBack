package main

import (
	"context"
	"log"
	"net"

	"eduData/bootstrap"
	pb "eduData/grpc"
	hrbustUg "eduData/school/hrbust/Ug"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// server 是 AuthService 的实现
type server struct {
	pb.UnimplementedAuthServiceServer
}

// Signin 实现
func (s *server) Signin(ctx context.Context, req *pb.SigninRequest) (*pb.SigninResponse, error) {
	logrus.Debugf("Signin request: %v", req)

	cookiej, err := hrbustUg.Signin(req.Username, req.Password)
	if err != nil {
		return &pb.SigninResponse{
			CookieJar:    []byte(""),
			ErrorMessage: err.Error(),
			Success:      false,
		}, nil
	}

	scj, err := pb.SerializeCookieJar(cookiej)
	if err != nil {
		return &pb.SigninResponse{
			CookieJar:    []byte(""),
			ErrorMessage: err.Error(),
			Success:      false,
		}, nil
	}

	return &pb.SigninResponse{
		CookieJar:    scj,
		ErrorMessage: "",
		Success:      true,
	}, nil
}

// GetData 实现
func (s *server) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	logrus.Debugf("GetData request: %v", req)

	cookiej, err := pb.DeserializeCookieJar(req.CookieJar)
	if err != nil {
		return &pb.GetDataResponse{
			Data:         []byte(""),
			ErrorMessage: err.Error(),
			Success:      false,
		}, nil
	}

	cData, err := hrbustUg.GetCourseByTime(cookiej, req.Year, req.Term)
	if err != nil {
		return &pb.GetDataResponse{
			Data:         []byte(""),
			ErrorMessage: err.Error(),
			Success:      false,
		}, nil
	}

	return &pb.GetDataResponse{
		Data:         *cData,
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func main() {
	bootstrap.Loadconfig()
	logrus.SetLevel(logrus.DebugLevel)
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
