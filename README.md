grpcintercept
===========

Generates boilerplate to help use the interceptor/middleware pattern with gRPC server service calls.
Desgined to work alongside the gRPC protobuf stub generation system, although this is not required.

#### Installation
Ensure Go is installed on your computer.
Run the standard go get as so:

	go get github.com/willtrking/grpcintercept


#### Usage
Takes advantage of the execellent `go generate` tool.

In each file you wish to generate boilerplate for, add the following comment

```go
//go:generate grpcintercept -Service SERVICENAME -GRPCRegister REGISTERCALL $GOFILE
```

SERVICENAME must be the name of the service you want to register with gRPC.
This should have all of the necessary functions attached to it, as described in your
gRPC protobuf stub file.


REGISTERCALL must be the server registration call generated in your gRPC protobuf stub file.
This should include the package in dot format, note that that package MUST be used in the file
your with your go:generate comment.


For example, if I have a `AccountManagementServer` in my stub file, which has a `RegisterAccountManagementServer` function, I would import that stub file with something like
`pb example/account/protobuf` and my REGISTERCALL would be `pb.RegisterAccountManagementServer`




Each generated file will have a `RegisterService` function which you should use in your server registration.




See the examples folder for more details
