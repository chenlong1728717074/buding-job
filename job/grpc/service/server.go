package service

import (
	"buding-job/job/core"
	"buding-job/job/grpc/to"
	"buding-job/job/handle"
	"buding-job/orm"
	"buding-job/orm/do"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	to.UnimplementedServerServer
}

func (*Server) Register(ctx context.Context, req *to.RegisterRequest) (*emptypb.Empty, error) {
	var m do.JobManagementDo
	orm.DB.First(&m, req.GetJobManagerId())
	if m.Id == 0 {
		return nil, errors.New("JobManagement NOT FOUND")
	}
	handle.JobManagerProcessor.RegisterInstance(core.NewInstance(req.ServiceAddr, req.JobManagerId, req.JobManagerName))
	return &emptypb.Empty{}, nil
}

func (*Server) Logout(context.Context, *to.RegisterRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
