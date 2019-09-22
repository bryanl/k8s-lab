package commands

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bryanl/k8s-lab/pkg/k8slab"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Run shell",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		runShell(ctx, clusterName, args)
	},
}

func runShell(ctx context.Context, clusterName string, args []string) {
	if err := k8slab.Shell(ctx, clusterName, args); err != nil {
		logrus.WithError(err).Error("run k8s-lab shell")
		os.Exit(1)
	}
}
