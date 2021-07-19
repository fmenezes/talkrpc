package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/fmenezes/talkrpc/grpc/service"
)

func stdinMode(cnn *grpc.ClientConn) error {
	client := pb.NewTalkRPCClient(cnn)
	stream, err := client.DoSomeWorkStream(context.Background())
	if err != nil {
		return err
	}
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("failed: %v", err)
				return
			}
			fmt.Println(res.Message)
		}
	}()

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		req := &pb.RequestResponse{Message: scan.Text()}
		fmt.Printf("Sending: %s\n", req.Message)
		err := stream.Send(req)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		if req.Message == "-" {
			err := stdinMode(cnn)
			if err != nil {
				log.Printf("failed: %s", err)
				return
			}
		}
	}

	client := pb.NewTalkRPCClient(cnn)
	res, err := client.DoSomeWork(context.Background(), req)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	fmt.Println(res.Message)
}
