package container

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Options struct {
	Mounts []mount.Mount
	Links  []string
}

func Run(ctx context.Context, options Options, imageName string, args ...string) error {
	canonicalImageName := "docker.io/" + imageName
	if !strings.Contains(imageName, "/") {
		canonicalImageName = "docker.io/library/" + imageName
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageSummaries, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return fmt.Errorf("list images: %w", err)
	}
	spew.Dump(imageSummaries)

	_, err = cli.ImagePull(ctx, canonicalImageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	containerConfig := &container.Config{
		Image: imageName,
		Cmd:   args,
	}
	hostConfig := &container.HostConfig{
		Links:       options.Links,
		Mounts:      options.Mounts,
		NetworkMode: "host",
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, out); err != nil {
		panic(err)
	}

	return nil
}
