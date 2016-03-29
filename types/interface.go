package types


type InterceptorData interface {
  Close() error
}

type Interceptor interface {
  Init() (InterceptorData, error)
}
