package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func main() {
	volumeName := "k8s-lab"

	// ensure volume exists
	if err := exec.Command("docker", "volume", "inspect", volumeName).Run(); err != nil {
		logger := logrus.WithField("volume-name", volumeName)
		logger.Info("creating docker volume")
		if err := exec.Command("docker", "volume", "create", volumeName).Run(); err != nil {
			logger.WithError(err).Error("unable to create docker volume")
			os.Exit(1)
		}
	}

	volumeMountArg := fmt.Sprintf("source=%s,target=/root/.kube", volumeName)

	execArgs := []string{
		"run",
		"-it",
		"--rm",
		"--net", "host",
		"-v", "//var/run/docker.sock:/var/run/docker.sock",
		"--mount", volumeMountArg,
		"bryanl/k8s-lab-wrapper",
	}
	execArgs = append(execArgs, os.Args[1:]...)

	cmd := exec.Command("docker", execArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Errorf("running docker")
	}
}
