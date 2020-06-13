package container

import (
	"context"
	"net/http"
	"testing"

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
		err = c.PullImage(image)
		So(err, ShouldBeNil)
		id, portBind, err := c.Create(image+":latest", []string{"ENV=test"}, map[string]string{"app": "test"}, []string{"8000/tcp"})
		So(err, ShouldBeNil)
		So(portBind, ShouldHaveLength, 1)
		_, err = http.Get("localhost:" + portBind["8000/tcp"] + "/ip")
		err = c.Destroy(id)
		So(err, ShouldBeNil)
	})

}
