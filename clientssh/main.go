package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"

	"github.com/fmenezes/talkrpc/common"
	"golang.org/x/crypto/ssh"
)

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		log.Fatalf("usage: %s <path> [message]", os.Args[0])
	}

	path := os.Args[1]

	key, err := ioutil.ReadFile(os.Getenv("SSH_KEY"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: os.Getenv("REMOTE_USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ssh_client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", os.Getenv("REMOTE_HOST")), config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer ssh_client.Close()

	unix_socket, err := ssh_client.Dial("unix", path)
	if err != nil {
		log.Fatalf("Failed to remote dial: %s", err)
	}
	defer unix_socket.Close()

	client := rpc.NewClient(unix_socket)
	defer client.Close()

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
