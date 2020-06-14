package container

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/docker/docker/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContainer(t *testing.T) {
	Convey("test container basic operation", t, func() {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		So(err, ShouldBeNil)
		c := NewClient(ctx, cli)
		image := "registry.cn-hangzhou.aliyuncs.com/lvanneo/httpbin"
		exposedPorts, portBindings, portBind, err := c.PrePareNetwork([]string{"8000/tcp"})
		// in the specific working, we will add new env here about port and ip
		id, err := c.Create(image+":latest", []string{"ENV=test"}, map[string]string{"app": "test"}, exposedPorts, portBindings)
		So(portBind, ShouldHaveLength, 1)
		fmt.Println(portBind["8000/tcp"])
		time.Sleep(2 * time.Second)
		res, err := http.Get("http://localhost:" + portBind["8000/tcp"] + "/ip")
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusOK)
		err = c.Destroy(id)
		So(err, ShouldBeNil)
	})

}
