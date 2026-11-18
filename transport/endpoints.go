package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/student-service/student"
)

// DTO для запросов и ответов
type CreateRequest struct {
	Student student.Student
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Student student.Student `json:"student"`
	Err     error           `json:"error,omitempty"`
}

type ChangeGroupRequest struct {
	ID    string `json:"id"`
	Group string `json:"group"`
}

type ChangeGroupResponse struct {
	Err error `json:"error,omitempty"`
}

// Endpoints – набор эндпоинтов сервиса студентов
type Endpoints struct {
	Create      endpoint.Endpoint
	GetByID     endpoint.Endpoint
	ChangeGroup endpoint.Endpoint
}

// MakeEndpoints создаёт эндпоинты для сервиса студентов
func MakeEndpoints(s student.Service) Endpoints {
	return Endpoints{
		Create:      makeCreateEndpoint(s),
		GetByID:     makeGetByIDEndpoint(s),
		ChangeGroup: makeChangeGroupEndpoint(s),
	}
}

func makeCreateEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Student)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		st, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Student: st, Err: err}, nil
	}
}

func makeChangeGroupEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeGroupRequest)
		err := s.ChangeGroup(ctx, req.ID, req.Group)
		return ChangeGroupResponse{Err: err}, nil
	}
}
