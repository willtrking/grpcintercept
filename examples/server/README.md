Server Example
===========

Full featured sample server which attaches a `github.com/gocraft/dbr` session to each gRPC service call

Run with
`go run main.go` inside the `server` folder

gRPC stubs created with

`protoc --go_out=plugins=grpc:. account.proto`


grpcintercept boilerplate created with

`go generate`

in the `server/account` folder
