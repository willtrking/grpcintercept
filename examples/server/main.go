package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	account "github.com/willtrking/grpcintercept/examples/server/account"
	db "github.com/willtrking/grpcintercept/examples/server/database"
	icpt "github.com/willtrking/grpcintercept/examples/server/intercept"
)

var (
	address = flag.String("address", "0.0.0.0", "The server address")
	port    = flag.Int("port", 10000, "The server port")
)

func main() {

	//Parse our flags
	flag.Parse()

	//Setup our database connection
	conn, cerr := db.GetConnection()

	if cerr != nil {
		panic(cerr)
	}

	//Setup our interceptor generation object
	interceptor := new(icpt.InitInterceptor)
	//Attach database connection to it
	interceptor.Conn = conn

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *address, *port))

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	grpclog.Println(" - LISTENING ON", fmt.Sprintf("%s:%d", *address, *port))

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	grpclog.Println(" - ATTACHING SERVICES")

	//Call the generated RegisterService call with our gRPC server
	//and our interceptor generation object
	account.RegisterService(grpcServer, interceptor)
	grpclog.Println("  - ATTACHED AccountService Service")
	grpclog.Println("SERVING")

	grpcServer.Serve(lis)

}
