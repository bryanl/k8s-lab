package commands

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bryanl/k8s-lab/pkg/k8slab"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize lab environment",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		runInit(ctx)
	},
}

func runInit(_ context.Context) {
	if err := k8slab.Init(clusterName); err != nil {
		logrus.WithError(err).Error("init k8s-lab")
		os.Exit(1)
	}
}
