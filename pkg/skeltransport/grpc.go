package skeltransport

import (
	"context"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/red010b37/sketetonservice/pb"
	"github.com/red010b37/sketetonservice/pkg/skelendpoint"

	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	login grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC GreeterServer.
func NewGRPCServer(endpoints svcendpoint.Endpoints, logger log.Logger) pb.LoginServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		login: grpctransport.NewServer(
			endpoints.LoginEndPoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
			options...,
		),
	}
}

// Login implementation of the method of the LoginService interface.
func (s *grpcServer) Login(ctx oldcontext.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, res, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.LoginResponse), nil
}

// decodeGRPCLoginRequest is a transport/grpc.DecodeRequestFunc that converts
// a gRPC login request to a user-domain login request.
func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return svcendpoint.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

// encodeGRPCLoginResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain login response to a gRPC login response.
func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(svcendpoint.LoginResponse)
	return &pb.LoginResponse{Token: res.Token}, nil
}
