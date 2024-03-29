package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/bryanl/k8s-lab/pkg/wrapper"
)

func main() {
	userArgs := os.Args[1:]

	command := ""
	if len(userArgs) > 0 {
		command = userArgs[0]
	}

	switch command {
	case "octant":
		runOctant()
	default:
		run(userArgs...)
	}
}

func runOctant() {
	or := wrapper.NewOctantRunner()
	octantBin, err := or.BinPath()
	if err != nil {
		logrus.WithError(err).Errorf("find octant binary")
		return
	}

	var buf bytes.Buffer

	r := newRunner()
	r.stdout = &buf

	if err := r.run("config"); err != nil {
		logrus.WithError(err).Errorf("retrieving config")
		return
	}

	configFile, err := ioutil.TempFile("", "kube-config")
	if err != nil {
		logrus.WithError(err).Errorf("opening temp config file")
		return
	}

	defer func() {
		if err := os.Remove(configFile.Name()); err != nil {
			logrus.WithError(err).Errorf("removing temp config file")
			return
		}
	}()

	if _, err := configFile.Write(buf.Bytes()); err != nil {
		logrus.WithError(err).Errorf("write temp config file")
		return
	}
	if err := configFile.Close(); err != nil {
		logrus.WithError(err).Errorf("closing temp config file")
		return
	}

	env := os.Environ()
	env = append(env,
		"OCTANT_ENABLE_APPLICATIONS=1",
		"OCTANT_DISABLE_OPEN_BROWSER=1",
	)

	cmd := exec.Command(octantBin, "--kubeconfig", configFile.Name())
	cmd.Env = env
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		logrus.WithError(err).Errorf("running Octant")
	}
}

func run(userArgs ...string) {
	r := newRunner()
	if err := r.run(userArgs...); err != nil {
		logrus.WithError(err).Errorf("running docker")
	}
}

type runner struct {
	stderr io.Writer
	stdout io.Writer
	stdin  io.Reader
}

func newRunner() *runner {
	return &runner{
		stderr: os.Stderr,
		stdout: os.Stdout,
		stdin:  os.Stdin,
	}
}

func (r *runner) run(userArgs ...string) error {
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
		// TODO: get this image from somewhere else
		"bryanl/k8s-lab-wrapper:0.1.0",
	}
	execArgs = append(execArgs, userArgs...)

	cmd := exec.Command("docker", execArgs...)
	cmd.Stderr = r.stderr
	cmd.Stdin = r.stdin
	cmd.Stdout = r.stdout

	return cmd.Run()
}
