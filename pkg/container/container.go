package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Volume struct {
	Source string
	Target string
}

type Mount struct {
	Source string
	Target string
}

type Options struct {
	Mounts  []Mount
	Volumes []Volume
	Links   []string
}

func Interactive(ctx context.Context, options Options, imageName string, args ...string) error {
	execArgs := []string{
		"run",
		"-it",
		"--rm",
		"--net", "host",
		"--add-host", "app.local:127.0.0.1",
	}
	for _, v := range options.Volumes {
		spec := fmt.Sprintf("%s:%s", v.Source, v.Target)
		execArgs = append(execArgs, "-v", spec)
	}
	for _, m := range options.Mounts {
		spec := fmt.Sprintf("source=%s,target=%s", m.Source, m.Target)
		execArgs = append(execArgs, "--mount", spec)
	}
	execArgs = append(execArgs, imageName)
	execArgs = append(execArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", execArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
