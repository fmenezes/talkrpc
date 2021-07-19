package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/fmenezes/talkrpc/rpc/common"
)

type App struct{}

func (a *App) DoSomeWork(req common.Request, res *common.Response) (err error) {
	fmt.Printf("Received: %s\n", req.Message)
	res.Message = fmt.Sprintf("Responded: %s", req.Message)
	return
}

func ServeAt(path string) (err error) {
	rpc.Register(&App{})

	listener, err := net.Listen("unix", path)
	if err != nil {
		return fmt.Errorf("unable to listen at %s: %s", path, err)
	}

	go rpc.Accept(listener)
	return
}

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		log.Printf("usage: %s <path>", os.Args[0])
		return
	}

	path := os.Args[1]

	err := ServeAt(path)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	defer os.Remove(path)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP)
	log.Printf("Received signal %s. Exiting...\n", <-signals)
}
