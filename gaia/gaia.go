package gaia

import (
	"context"
	"github.com/ThunderYurts/Gaia/container"
	"github.com/ThunderYurts/Gaia/gserver"
	"github.com/ThunderYurts/Gaia/monitor"
	"github.com/docker/docker/client"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

type Gaia struct {
	server     *gserver.Server
	monitor    *monitor.SimpleMonitor
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         *sync.WaitGroup
	ip         string
	port       string
	name       string
}

// NewGaia is a help function
func NewGaia(ip string, port string, name string, zkAddr []string) Gaia {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	conn, _, err := zk.Connect(zkAddr, 5*time.Second)
	if err != nil {
		panic(err)
	}
	containerCli := container.NewClient(ctx, cli)
	s := gserver.NewServer(ctx, &wg, &containerCli, ip)
	mon := monitor.NewSimpleMonitor(ctx, cli, conn, name, ip+port)
	return Gaia{
		server:     &s,
		monitor:    &mon,
		ctx:        ctx,
		cancelFunc: cancel,
		wg:         &wg,
		ip:         ip,
		port:       port,
		name:       name,
	}
}

func (g *Gaia) Start() {
	err := g.server.Start(g.port)
	if err != nil {
		panic(err)
	}
	go func() {
		err := g.monitor.SyncNodeStat()
		if err != nil {
			panic(err)
		}
	}()
}

func (g *Gaia) Stop() {
	g.cancelFunc()
	g.wg.Wait()
}
