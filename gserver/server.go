package gserver

import (
	"context"
	"fmt"
	"net"
	sync "sync"

	"github.com/ThunderYurts/Gaia/container"
	"github.com/ThunderYurts/Gaia/gconst"
	mapset "github.com/deckarep/golang-set"
	"google.golang.org/grpc"
)

// Server will start for grpc
type Server struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	client   *container.Client
	ip       string
	yurtpool mapset.Set
}

// NewServer is a help function
func NewServer(ctx context.Context, wg *sync.WaitGroup, client *container.Client, ip string) Server {
	return Server{
		ctx:      ctx,
		wg:       wg,
		client:   client,
		ip:       ip,
		yurtpool: mapset.NewSet(),
	}
}

// Delete put key value pair into storage
func (s *Server) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
	// create yurt to yurt pool
	exposedPorts, portBindings, portBind, err := s.client.PrePareNetwork([]string{gconst.ActionPort, gconst.SyncPort})
	if err != nil {
		panic(err)
		return &CreateReply{Code: CreateCode_CREATE_ERROR}, nil
	}

	env := []string{"HOST_IP=" + s.ip, "ACTION_PORT=:" + portBind[gconst.ActionPort], "SYNC_PORT=:" + portBind[gconst.SyncPort], "SERVICE_NAME:" + in.ServiceName}
	fmt.Printf("env: %v\n", env)
	_, err = s.client.Create(gconst.YurtImage, env, map[string]string{"app": "thunderyurt"}, exposedPorts, portBindings)

	if err != nil {
		panic(err)
		return &CreateReply{Code: CreateCode_CREATE_ERROR}, nil
	}
	return &CreateReply{Code: CreateCode_CREATE_SUCCESS}, nil
}

// Start will run a server
func (s *Server) Start(port string) error {
	gaiaServer := grpc.NewServer()
	RegisterBreedServer(gaiaServer, s)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		fmt.Printf("action server listen on %s\n", port)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			select {
			case <-s.ctx.Done():
				{
					fmt.Println("action get Done")
					gaiaServer.GracefulStop()
					return
				}
			}
		}(s.wg)
		gaiaServer.Serve(lis)
	}(s.wg)
	return nil

}
