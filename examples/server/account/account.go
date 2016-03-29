package account

//go:generate grpcintercept -Service AccountService -GRPCRegister pb.RegisterAccountManagementServer $GOFILE

import (
	icpt "github.com/willtrking/grpcintercept/examples/server/intercept"
	"golang.org/x/net/context"

	pb "github.com/willtrking/grpcintercept/examples/server/account/protobuf"
)

type AccountService struct{}

func (a *AccountService) GetAccount(ctx context.Context, account *pb.Account, idat *icpt.InterceptorStore) (*pb.Account, error) {

	//Here I could use my idat.Db dbr Session to interact with my database

	test := pb.Account{1, "GET William King", "test@example.com"}

	return &test, nil
}

func (a *AccountService) CreateAccount(ctx context.Context, account *pb.Account, idat *icpt.InterceptorStore) (*pb.Account, error) {

	//Here I could use my idat.Db dbr Session to interact with my database

	test := pb.Account{2, "CREATE William King", "test@example.com"}

	return &test, nil
}
