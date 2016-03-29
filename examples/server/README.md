Server Example
===========

Full featured sample server which attaches a postgres `github.com/gocraft/dbr` session to each gRPC service call.

Postgres through `github.com/lib/pq`

Run with
`go run main.go`

gRPC stubs created with

`protoc --go_out=plugins=grpc:. account.proto`


grpcintercept boilerplate created with

`go generate`

in the `server/account` folder
