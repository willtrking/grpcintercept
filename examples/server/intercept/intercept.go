package intercept

import (
  "github.com/gocraft/dbr"
  grpci "github.com/willtrking/grpcintercept/types"
)

//We want a dbr session in each of our gRPC calls
//This is where we do that setup

type InterceptorStore struct {
  Db *dbr.Session
}

type InitInterceptor struct {
  Conn *dbr.Connection
}

func (i *InitInterceptor) Init() (grpci.InterceptorData, error) {
  m := new(InterceptorStore)
  m.Db = i.Conn.NewSession(nil)
  return m, nil
}

func (is InterceptorStore) Close() error {
  is.Db = nil
  return nil
}
