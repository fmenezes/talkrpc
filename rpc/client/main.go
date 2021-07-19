package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/fmenezes/talkrpc/rpc/common"
)

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		log.Fatalf("usage: %s <path> [message]", os.Args[0])
	}

	path := os.Args[1]

	cnn, err := net.Dial("unix", path)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}
	client := rpc.NewClient(cnn)

	req := &common.Request{Message: "No message"}

	if len(os.Args) > 2 {
		req.Message = os.Args[2]
	}

	var res common.Response
	err = client.Call("App.DoSomeWork", req, &res)
	if err != nil {
		log.Printf("error in rpc: %s", err)
		return
	}
	fmt.Println(res.Message)
}
