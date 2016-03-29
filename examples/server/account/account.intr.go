// Code generated by grpcintercept
// source: account.go
// DO NOT EDIT!

package account

import (
	grpcintercept "github.com/willtrking/grpcintercept/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/willtrking/grpcintercept/examples/server/account/protobuf"

	icpt "github.com/willtrking/grpcintercept/examples/server/intercept"
)

type AccountServiceInterceptor struct {
	service *AccountService
	i       grpcintercept.Interceptor
}

func RegisterService(s *grpc.Server, i grpcintercept.Interceptor) {
	srv := new(AccountServiceInterceptor)
	srv.i = i
	pb.RegisterAccountManagementServer(s, srv)
}

func (a *AccountServiceInterceptor) CreateAccount(ctx context.Context, account *pb.Account) (*pb.Account, error) {
	di, _ := a.i.Init()

	defer func(di grpcintercept.InterceptorData) {

		ce := di.Close()
		if ce != nil {
			grpclog.Println("Failed to close InterceptorData on CreateAccount ", ce)
		}

	}(di)

	return a.service.CreateAccount(ctx, account, di.(*icpt.InterceptorStore))
}

func (a *AccountServiceInterceptor) GetAccount(ctx context.Context, account *pb.Account) (*pb.Account, error) {
	di, _ := a.i.Init()

	defer func(di grpcintercept.InterceptorData) {

		ce := di.Close()
		if ce != nil {
			grpclog.Println("Failed to close InterceptorData on GetAccount ", ce)
		}

	}(di)

	return a.service.GetAccount(ctx, account, di.(*icpt.InterceptorStore))
}
