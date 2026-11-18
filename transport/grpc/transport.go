package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"example.com/student-service/pb"
	"example.com/student-service/student"
	"example.com/student-service/transport"
)

// grpcServer реализует pb.StudentServiceServer, используя go-kit grpc-хендлеры
type grpcServer struct {
	pb.UnimplementedStudentServiceServer

	create      grpctransport.Handler
	getByID     grpctransport.Handler
	changeGroup grpctransport.Handler
}

// NewGRPCServer создаёт gRPC-сервер на основе эндпоинтов Go kit
func NewGRPCServer(endpoints transport.Endpoints) pb.StudentServiceServer {
	options := []grpctransport.ServerOption{}

	return &grpcServer{
		create: grpctransport.NewServer(
			endpoints.Create,
			decodeCreateRequest,
			encodeCreateResponse,
			options...,
		),
		getByID: grpctransport.NewServer(
			endpoints.GetByID,
			decodeGetByIDRequest,
			encodeGetByIDResponse,
			options...,
		),
		changeGroup: grpctransport.NewServer(
			endpoints.ChangeGroup,
			decodeChangeGroupRequest,
			encodeChangeGroupResponse,
			options...,
		),
	}
}

// --- Реализация серверного интерфейса protobuf ---

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) ChangeGroup(ctx context.Context, req *pb.ChangeGroupRequest) (*pb.ChangeGroupResponse, error) {
	_, resp, err := s.changeGroup.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ChangeGroupResponse), nil
}

// --- Декодеры: protobuf -> DTO Go kit ---

func decodeCreateRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.CreateRequest)
	return transport.CreateRequest{
		Student: student.Student{
			FullName: r.FullName,
			Group:    r.GroupName,
			GPA:      r.Gpa,
		},
	}, nil
}

func decodeGetByIDRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: r.Id}, nil
}

func decodeChangeGroupRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.ChangeGroupRequest)
	return transport.ChangeGroupRequest{ID: r.Id, Group: r.GroupName}, nil
}

// --- Энкодеры: DTO Go kit -> protobuf ---

func encodeCreateResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(transport.CreateResponse)
	return &pb.CreateResponse{Id: r.ID, Error: errorString(r.Err)}, nil
}

func encodeGetByIDResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(transport.GetByIDResponse)
	return &pb.GetByIDResponse{
		Student: &pb.Student{
			Id:         r.Student.ID,
			FullName:   r.Student.FullName,
			GroupName:  r.Student.Group,
			Gpa:        r.Student.GPA,
			EnrolledOn: r.Student.EnrolledOn,
		},
		Error: errorString(r.Err),
	}, nil
}

func encodeChangeGroupResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(transport.ChangeGroupResponse)
	return &pb.ChangeGroupResponse{Error: errorString(r.Err)}, nil
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
