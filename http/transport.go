package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log" // <-- ИСПРАВЛЕНО: убран /kit/
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/student-service/transport"
)

// NewService связывает эндпоинты Go kit с HTTP-транспортом
func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/students").Handler(kithttp.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/students/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PATCH").Path("/students/{id}/group").Handler(kithttp.NewServer(
		endpoints.ChangeGroup,
		decodeChangeGroupRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Student); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeChangeGroupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var body struct {
		Group string `json:"group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transport.ChangeGroupRequest{ID: vars["id"], Group: body.Group}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()}) // <-- ИСПРАВЛЕНО: добавлено _ =
}
