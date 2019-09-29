package commands

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bryanl/k8s-lab/pkg/k8slab"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show kube config",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		runConfig(ctx, clusterName)
	},
}

func runConfig(ctx context.Context, clusterName string) {
	if err := k8slab.Config(ctx, clusterName); err != nil {
		logrus.WithError(err).Error("run k8s-lab config")
		os.Exit(1)
	}
}
