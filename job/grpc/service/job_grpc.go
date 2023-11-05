package service

import (
	"buding-job/job/grpc/to"
	"buding-job/job/handle"
	"buding-job/orm"
	"buding-job/orm/do"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type JobService struct {
	to.UnimplementedJobServer
}

func NewJobService() *JobService {
	return &JobService{}
}

func (*JobService) Callback(ctx context.Context, resp *to.CallbackResponse) (*emptypb.Empty, error) {
	var jobLog do.JobLogDo
	orm.DB.First(&jobLog, resp.GetId())
	if jobLog.Id == 0 {
		return nil, errors.New("entry does not exist")
	}
	//async handle log
	go handle.JobMonitor.Callback(&jobLog, resp)
	return &emptypb.Empty{}, nil
}
