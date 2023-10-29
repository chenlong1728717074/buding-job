package job

import (
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
}
func (app *JobApp) Stop() {
	app.grpcServer.Stop()
}
