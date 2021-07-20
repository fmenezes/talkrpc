package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
		return err
	}
	defer os.Remove(path)
	grpcServer := grpc.NewServer()
	pb.RegisterTalkRPCServer(grpcServer, &server{})

	sig := make(chan os.Signal)
	errChan := make(chan error)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)

	select {
	case signal := <-sig:
		log.Printf("Received signal %s. Exiting...\n", signal)
	case err := <-errChan:
		return err
	}

	return nil
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
