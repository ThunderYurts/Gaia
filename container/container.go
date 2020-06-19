package container

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	retry "gopkg.in/matryer/try.v1"
)

// Client will hold connection will docker
type Client struct {
	ctx     context.Context
	client  *client.Client
	portSet mapset.Set
}

// NewClient is a help function
func NewClient(ctx context.Context, client *client.Client) Client {
	return Client{
		ctx:     ctx,
		client:  client,
		portSet: mapset.NewSet(),
	}
}

func (c *Client) getFreePort() (int, error) {
	var addr *net.TCPAddr
	err := retry.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		randomPort := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000) + 40000)
		addr, err = net.ResolveTCPAddr("tcp", "localhost:"+randomPort)
		if err != nil {
			return
		}
		if c.portSet.Contains(addr.Port) {
			err = errors.New("port has pre dispatched")
		} else {
			err = nil
		}
		return
	})
	c.portSet.Add(addr.Port)

	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// PullImage will pull image
func (c *Client) PullImage(image string) error {
	_, err := c.client.ImagePull(c.ctx, image, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

// PrePareNetwork will
func (c *Client) PrePareNetwork(exportbind []string) (exposedPorts nat.PortSet, portBindings nat.PortMap, exports map[string]string, err error) {
	exposedPorts = nat.PortSet{}
	portBindings = nat.PortMap{}
	exports = make(map[string]string)
	err = nil
	for _, external := range exportbind {
		exposedPorts[nat.Port(external)] = struct{}{}
		port, portErr := c.getFreePort()
		if portErr != nil {
			panic(portErr)
			err = portErr
			return
		}
		p := strconv.Itoa(port)
		exports[external] = p
		portBindings[nat.Port(external)] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: p}}
	}
	return
}

func randomName(n int, allowedChars ...[]rune) string {
	var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// Create a new container and use net=host(v1), will return error about do not have enough port to deploy
// will pull image if it is not existed
// pay attention to the result, we can not guarantee that the server in container is running and
func (c *Client) Create(image string, env []string, labels map[string]string, exposedPorts nat.PortSet, portBindings nat.PortMap) (string, error) {

	err := c.PullImage(image)
	if err != nil {
		return "", err
	}
	// TODO add PORT ENV
	container, err := c.client.ContainerCreate(c.ctx, &container.Config{
		Image:        image,
		Env:          append(env, "YURT_NAME="+randomName(10)),
		Labels:       labels,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
	},
		nil,
		"",
	)
	if err != nil {
		fmt.Printf("117 err : %v\n", err)
		return "", err
	}
	if err := c.client.ContainerStart(c.ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		// clean this container
		_ = c.client.ContainerRemove(c.ctx, container.ID, types.ContainerRemoveOptions{Force: true})
		return "", err
	}

	return container.ID, nil
}

// Destroy will destroy container by name
func (c *Client) Destroy(ID string) error {
	err := c.client.ContainerRemove(c.ctx, ID, types.ContainerRemoveOptions{Force: true})
	return err
}
