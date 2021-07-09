# talkrpc

## How to use this?

Server (process that listens): `go run server/main.go ./test`

Client (process that connects): `go run client/main.go ./test "My Message"`

## How to use this over ssh?

Server (process that listens should be stated on remote machine): `go run server/main.go ./test`

Client SSH (process that connects should be started on host machine): `SSH_KEY=~/.ssh/id_rsa REMOTE_HOST=<machine> REMOTE_USER=<user> go run clientssh/main.go /absolute/path/to/test "My Message"`
