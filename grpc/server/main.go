package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/fmenezes/talkrpc/grpc/service"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTalkRPCServer
}

func (s *server) DoSomeWork(ctx context.Context, req *pb.RequestResponse) (*pb.RequestResponse, error) {
	fmt.Printf("Received: %s\n", req.Message)
	return &pb.RequestResponse{
		Message: fmt.Sprintf("Responded: %s", req.Message),
	}, nil
}
func (s *server) DoSomeWorkStream(stream pb.TalkRPC_DoSomeWorkStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received: %s\n", req.Message)
		err = stream.Send(&pb.RequestResponse{
			Message: fmt.Sprintf("Responded: %s", req.Message),
		})
		if err != nil {
			return err
		}
	}
}

func serveAt(path string) error {
	lis, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer os.Remove(path)
	grpcServer := grpc.NewServer()
	pb.RegisterTalkRPCServer(grpcServer, &server{})
	return grpcServer.Serve(lis)
}

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		log.Printf("usage: %s <path>", os.Args[0])
		return
	}

	path := os.Args[1]

	err := serveAt(path)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
}
