package k8slab

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/cluster"

	"github.com/bryanl/k8s-lab/pkg/container"
)

func Init(clusterName string) error {
	logger := logrus.WithField("cluster-name", clusterName)

	isKnown, err := cluster.IsKnown(clusterName)
	if err != nil {
		return fmt.Errorf("check if cluster %s is known: %w", clusterName, err)
	}

	if isKnown {
		logger.Info("cluster already exists")
		return nil
	}

	kindContext := cluster.NewContext(clusterName)

	logger.Info("init cluster")

	if err := kindContext.Create(); err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}

	p, err := configPath(clusterName)
	if err != nil {
		return fmt.Errorf("get cluster config path: %w", err)
	}
	if err := os.Chmod(p, 0644); err != nil {
		return fmt.Errorf("update cluster config permissions: %w", err)
	}

	return nil
}

func Delete(clusterName string) error {
	logger := logrus.WithField("cluster-name", clusterName)

	kindContext := cluster.NewContext(clusterName)

	logger.Info("delete cluster")

	if err := kindContext.Delete(); err != nil {
		return fmt.Errorf("delete cluster: %w", err)
	}

	return nil
}

func Shell(ctx context.Context, clusterName string, args []string) error {
	targetDir := "/home/k8slab/.kube"
	targetConfigName := "kind-config-k8s-lab"
	options := container.Options{}

	runningInDocker := checkFileExists("/.dockerenv")
	if runningInDocker {
		logrus.Debug("running in docker; using k8s-lab mount")
		options.Mounts = []container.Mount{
			{
				Source: "k8s-lab",
				Target: targetDir,
			},
		}
	} else {
		source, err := configPath(clusterName)
		if err != nil {
			return fmt.Errorf("get cluster config path: %w", err)
		}

		options.Volumes = []container.Volume{
			{
				Source: source,
				Target: filepath.Join(targetDir, targetConfigName),
			},
		}
	}

	err := container.Interactive(ctx, options, "bryanl/k8s-lab", args...)
	if err != nil {
		return fmt.Errorf("unable to run container: %w", err)
	}

	return nil
}

func Config(ctx context.Context, clusterName string) error {
	args := []string{
		"cat", ".kube/kind-config-k8s-lab",
	}
	return Shell(ctx, clusterName, args)
}

func configPath(clusterName string) (string, error) {
	dir, err := kubeDir()
	if err != nil {
		return "", err
	}

	configName := fmt.Sprintf("kind-config-%s", clusterName)
	configPath := filepath.Join(dir, configName)

	return configPath, nil
}

func kubeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to find current user: %w", err)
	}

	return filepath.Join(u.HomeDir, ".kube"), nil
}

func checkFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
