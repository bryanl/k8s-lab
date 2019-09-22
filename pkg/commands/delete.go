package commands

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bryanl/k8s-lab/pkg/k8slab"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete lab environment",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		runDelete(ctx)
	},
}

func runDelete(_ context.Context) {
	if err := k8slab.Delete(clusterName); err != nil {
		logrus.WithError(err).Error("delete k8s-lab")
		os.Exit(1)
	}
}
