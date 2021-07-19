# talkrpc

## How to use `net/rpc`?

Server (process that listens): `go run rpc/server/main.go ./test.sock`

Client (process that connects): `go run rpc/client/main.go ./test.sock "My Message"`

## How to use `net/rpc` over ssh?

Server (process that listens should be started on remote machine): `go run rpc/server/main.go ./test.sock`

Client SSH (process that connects to host machine): `SSH_KEY=~/.ssh/id_rsa REMOTE_HOST=<machine> REMOTE_USER=<user> go run rpc/clientssh/main.go /absolute/path/to/test.sock "My Message"`
