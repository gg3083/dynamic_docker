package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	round := rand.Int()
	//创建一个本机docker客户端
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	runPort := 9000
	docker1 := Docker{c}

	ports := make(nat.PortSet)
	port, err := nat.NewPort("tcp", fmt.Sprintf("%d", runPort))
	if err != nil {
		panic(err)
	}
	ports[port] = struct{}{}

	portMap := make(nat.PortMap)
	portMap[port] = []nat.PortBinding{{
		HostIP:   "0.0.0.0",
		HostPort: fmt.Sprintf("%d", runPort),
	}}

	if err = docker1.Run(
		context.Background(),
		fmt.Sprintf("qrcode_%d", round/1e9),
		"qrcode:v1",
		[]string{fmt.Sprintf("port=%d", runPort)},
		ports,
		&container.HostConfig{
			PortBindings: portMap,
		},
	); err != nil {
		panic(err)
	}

	fmt.Println("container started")

	////删除container
	//err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
	//	Force: true,
	//})
	//if err != nil {
	//	panic(err)
	//}

}

type Docker struct {
	client *client.Client
}

func (d Docker) Run(ctx context.Context, name string, image string, env []string, ports nat.PortSet, hostConfig *container.HostConfig) error {
	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: ports,
		Env:          env,
	},
		hostConfig, nil, nil, name)
	if err != nil {
		return err
	}
	if err = d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
