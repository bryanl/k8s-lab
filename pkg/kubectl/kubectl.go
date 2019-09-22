package kubectl

import (
	"context"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"

	"github.com/bryanl/k8s-lab/pkg/container"
)

func Run(ctx context.Context, clusterName string, args ...string) error {
	imageName := "lachlanevenson/k8s-kubectl:v1.15.4"

	kubeHome, err := kubeDir()
	if err != nil {
		return err
	}
	configName := fmt.Sprintf("kind-config-%s", clusterName)
	configPath := filepath.Join(kubeHome, configName)

	targetPath := "/kube/config"

	options := container.Options{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: configPath,
				Target: targetPath,
			},
		},
	}

	args = append([]string{"--kubeconfig", targetPath}, args...)

	return container.Run(ctx, options, imageName, args...)
}

func kubeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to find current user: %w", err)
	}

	return filepath.Join(u.HomeDir, ".kube"), nil
}
