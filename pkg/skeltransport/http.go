package skeltransport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/red010b37/sketetonservice/pkg/skelendpoint"
	"github.com/red010b37/sketetonservice/pkg/skelerror"
	"net/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints svcendpoint.Endpoints, logger log.Logger) http.Handler {
	m := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	// GET /health         retrieves service heath information
	// POST /login         allows people to login

	m.Methods("GET").Path("/_ah/health").Handler(httptransport.NewServer(
		endpoints.HealthEndpoint,
		DecodeHTTPHealthRequest,
		EncodeHTTPGenericResponse,
		options...,
	))

	m.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.LoginEndPoint,
		DecodeHTTPLoginRequest,
		EncodeHTTPGenericResponse,
		options...,
	))

	//m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//
	//	w.Write([]byte("FOO:" + authconfig.AppConfig.DBName))
	//
	//	//logger.Log("ov", os.Getenv("FOO"))
	//
	//})

	return m
}

// DecodeHTTPHealthRequest method.
func DecodeHTTPHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return svcendpoint.HealthRequest{}, nil
}

// DecodeHTTPGreetingRequest method.
func DecodeHTTPLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {

	// Check we have a body
	//if r.GetBody == nil {
	//	return nil, &autherror.NoBody
	//}

	// try to decode the body
	var req svcendpoint.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		return nil, &skelerror.JSONDecodeError
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))

	// cast it to an app error
	appErr := err.(*skelerror.AuthServiceError)
	wrappedErr := errorWrapper{Error: appErr.Error(), Key: appErr.Key()}

	json.NewEncoder(w).Encode(wrappedErr)
}

func err2code(err error) int {

	appErr := err.(*skelerror.AuthServiceError)

	switch appErr {
	case &skelerror.JSONDecodeError, &skelerror.NoBody:
		return http.StatusBadRequest
	case &skelerror.InvalidCredentials:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Key   string `json:"key,omitempty"`
	Error string `json:"error,omitempty"`
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(svcendpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
