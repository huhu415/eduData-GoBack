package main

import (
	"context"
	"log"
	"time"

	pb "eduData/grpc" // 替换为生成的 pb.go 和 _grpc.pb.go 文件的实际路径

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用 Signin 方法
	signinResp, err := c.Signin(ctx, &pb.SigninRequest{Username: "test", Password: "password"})
	if err != nil {
		log.Fatalf("could not signin: %v", err)
	}
	log.Printf("Signin response: %v", signinResp)

	// 调用 GetData 方法
	getDataResp, err := c.GetData(ctx, &pb.GetDataRequest{CookieJar: signinResp.CookieJar, ModuleId: "module1"})
	if err != nil {
		log.Fatalf("could not get data: %v", err)
	}
	log.Printf("GetData response: %v", getDataResp)
}
