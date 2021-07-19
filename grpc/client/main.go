package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/fmenezes/talkrpc/grpc/service"
)

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		log.Fatalf("usage: %s <path> [message]", os.Args[0])
	}

	path := os.Args[1]

	cnn, err := grpc.Dial(
		path,
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	defer cnn.Close()

	req := &pb.RequestResponse{Message: "No message"}
	if len(os.Args) > 2 {
		req.Message = os.Args[2]
	}

	client := pb.NewTalkRPCClient(cnn)
	res, err := client.DoSomeWork(context.Background(), req)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	fmt.Println(res.Message)
}
