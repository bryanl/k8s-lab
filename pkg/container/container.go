package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types/mount"
)

type Options struct {
	Mounts []mount.Mount
	Links  []string
}

func Interactive(ctx context.Context, options Options, imageName string, args ...string) error {
	execArgs := []string{"run", "-it", "--rm", "--net", "host"}
	for _, m := range options.Mounts {
		spec := fmt.Sprintf("%s:%s", m.Source, m.Target)
		execArgs = append(execArgs, "-v", spec)
	}
	execArgs = append(execArgs, imageName)
	execArgs = append(execArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", execArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
