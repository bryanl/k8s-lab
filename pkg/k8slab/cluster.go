package k8slab

import (
	"context"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
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
	logger := logrus.WithField("cluster-name", clusterName)

	logger.Info("run shell")

	kubeHome, err := kubeDir()
	if err != nil {
		return err
	}

	configName := fmt.Sprintf("kind-config-%s", clusterName)
	configPath := filepath.Join(kubeHome, configName)

	targetPath := "/home/k8slab/.kube/config"

	options := container.Options{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: configPath,
				Target: targetPath,
			},
		},
	}

	err = container.Interactive(ctx, options, "bryanl/k8slab", args...)
	if err != nil {
		return fmt.Errorf("unable to run container: %w", err)
	}

	return nil
}

func kubeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to find current user: %w", err)
	}

	return filepath.Join(u.HomeDir, ".kube"), nil
}
