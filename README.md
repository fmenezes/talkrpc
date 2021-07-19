# talkrpc

## How to use `net/rpc`?

Server (process that listens): `go run rpc/server/main.go ./test.sock`

Client (process that connects): `go run rpc/client/main.go ./test.sock "My Message"`

## How to use `net/rpc` over ssh?

Server (process that listens should be started on remote machine): `go run rpc/server/main.go ./test.sock`

Client SSH (process that connects to host machine): `SSH_KEY=~/.ssh/id_rsa REMOTE_HOST=<machine> REMOTE_USER=<user> go run rpc/clientssh/main.go /absolute/path/to/test.sock "My Message"`

## How to use `grpc`?

Build: `cd service && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service.proto`

Server (process that listens): `go run grpc/server/main.go ./test.sock`

Client (process that connects): `go run grpc/client/main.go ./test.sock "My Message"`
for stream/stdin mode: `go run grpc/client/main.go ./test.sock -`
