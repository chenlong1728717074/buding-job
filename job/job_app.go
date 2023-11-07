package job

import (
	"buding-job/job/grpc/service"
	"buding-job/job/grpc/to"
	"buding-job/job/handle"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type JobApp struct {
	grpcServer *grpc.Server
}

func NewJobApp() *JobApp {
	return &JobApp{
		grpcServer: grpc.NewServer(),
	}
}

func (app *JobApp) Start() {
	handle.JobManagerProcessor.Start()
	handle.JobSchedule.Start()
	handle.JobMonitor.Start()
	app.registerServer()
	app.startGrpc()
}

func (app *JobApp) startGrpc() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 8082))
	if err != nil {
		panic(err)
	}
	go func() {
		err := app.grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Grpc Service startup failed:%s", err.Error())
		}
	}()
	log.Println("Grpc server startup...")
}
func (app *JobApp) Stop() {
	app.grpcServer.Stop()
}

func (app *JobApp) registerServer() {
	to.RegisterServerServer(app.grpcServer, service.NewServerService())
	to.RegisterJobServer(app.grpcServer, service.NewJobService())
}
