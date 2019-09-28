package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	clusterName string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-lab",
	Short: "k8s lab/workshop tool",
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&clusterName, "cluster-name", "", "k8s-lab", "Cluster name")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(shellCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
