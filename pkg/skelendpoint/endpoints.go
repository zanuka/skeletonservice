package svcendpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/red010b37/sketetonservice/pkg/skelsvc"
)

// Endpoints collects all of the endpoints that compose a auth service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	HealthEndpoint endpoint.Endpoint // used by app engine
	LoginEndPoint  endpoint.Endpoint
}

// MakeServerEndpoints returns service Endoints, and wires in all the provided middlewares.
func MakeServerEndpoints(s skelsvc.Service, logger log.Logger) Endpoints {

	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
		healthEndpoint = LoggingMiddleware(log.With(logger, "method", "Health"))(healthEndpoint)
	}

	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(s)
		loginEndpoint = LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
	}

	return Endpoints{
		HealthEndpoint: healthEndpoint,
		LoginEndPoint:  loginEndpoint,
	}

}

// MakeHealthEndpoint constructs a Health endpoint wrapping the service.
func MakeHealthEndpoint(s skelsvc.Service) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health()
		return HealthResponse{Healthy: healthy}, nil
	}

}

// MakeLoginEnpoint constructs a login endpoint wrapping the service.
func MakeLoginEndpoint(s skelsvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(LoginRequest)
		token, err := s.Login(req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return LoginResponse{Token: token}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// HEALTH REQUEST
//------------------------------------------------------------------------------------
// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// Failed implements Failer.
func (r HealthResponse) Failed() error { return r.Err }

// LOGN REQUEST
//------------------------------------------------------------------------------------
// LoginRequest collects the request parameters for the Login method.
type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// LoginResponse collects the response values for the Login method.
type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}

// Failed implements Failer.
func (r LoginResponse) Failed() error { return r.Err }
