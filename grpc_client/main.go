package main

import (
	"context"
	"time"

	"eduData/bootstrap"
	pb "eduData/grpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	bootstrap.Loadconfig()

	logrus.SetLevel(logrus.DebugLevel)
	conn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 调用 Signin 方法
	signinResp, err := c.Signin(ctx, &pb.SigninRequest{Username: "2204010417", Password: "13737826060a@Hlg10214"})
	if err != nil {
		logrus.Errorf("could not signin: %v", err)
	}
	if !signinResp.Success {
		logrus.Errorf("could not signin: %v", signinResp.ErrorMessage)
	}
	logrus.Debugf("Signin response: %v", signinResp)

	// 调用 GetData 方法
	getDataResp, err := c.GetData(ctx, &pb.GetDataRequest{CookieJar: signinResp.CookieJar, Year: "44", Term: "2"})
	if err != nil {
		logrus.Errorf("could not get data: %v", err)
	}
	logrus.Debugf("GetData response: %v", getDataResp)
}
