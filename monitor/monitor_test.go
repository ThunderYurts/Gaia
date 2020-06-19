package monitor

import (
	"context"
	"fmt"
	"github.com/ThunderYurts/Gaia/gconst"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"

	"github.com/docker/docker/client"
)

func TestEverything(t *testing.T) {
	// Maintaining this map is error-prone and cumbersome (note the subtle bug):
	fs := map[string]func(*testing.T){
		"testMonitor": testMonitor,
	}
	// You may be able to use the `importer` package to enumerate tests instead,
	// but that starts getting complicated.
	for name, f := range fs {
		ids, err := setup()
		if err != nil {
			t.Error(err.Error())
		}
		t.Run(name, f)
		teardown(ids)
	}
}
func teardown(ids []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	for _, c := range ids {
		err := cli.ContainerRemove(ctx, c, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}
}

func setup() ([]string, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	var ids []string
	for i := 0; i < 4; i = i + 1 {
		c, err := cli.ContainerCreate(ctx, &container.Config{Image: "registry.cn-hangzhou.aliyuncs.com/lvanneo/httpbin", Labels: map[string]string{"app": "thunderyurt"}}, nil, nil, "")
		if err != nil {
			panic(err)
		}
		ids = append(ids, c.ID)
	}

	for _, c := range ids {
		err := cli.ContainerStart(ctx, c, types.ContainerStartOptions{})
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func testMonitor(t *testing.T) {
	convey.Convey("test monitor", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		cli, err := client.NewEnvClient()
		convey.So(err, convey.ShouldBeNil)
		time.Sleep(2 * time.Second)

		conn, _, err := zk.Connect([]string{"localhost:2181"}, 5*time.Second)
		convey.So(err, convey.ShouldBeNil)
		exist, _, err := conn.Exists(gconst.GaiaRoot)
		convey.So(err, convey.ShouldBeNil)
		if !exist {
			_, err = conn.Create(gconst.GaiaRoot, []byte{}, 0, zk.WorldACL(zk.PermAll))
		}
		convey.So(err, convey.ShouldBeNil)
		monitor := NewSimpleMonitor(ctx, cli, conn, "test", "127.0.0.1:30000")
		go func() {
			err = monitor.SyncNodeStat()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
		time.Sleep(5 * time.Second)
		convey.So(err, convey.ShouldBeNil)
		_, stat, err := conn.Get(gconst.GaiaRoot + "/test")
		convey.So(err, convey.ShouldBeNil)
		versionOne := stat.Version
		time.Sleep(5 * time.Second)
		_, stat, err = conn.Get(gconst.GaiaRoot + "/test")
		convey.So(err, convey.ShouldBeNil)
		versionTwo := stat.Version
		convey.So(versionOne, convey.ShouldBeLessThan, versionTwo)
		cancel()
		err = cli.Close()
		convey.So(err, convey.ShouldBeNil)
		conn.Close()
	})
}
