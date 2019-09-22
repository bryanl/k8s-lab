package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bryanl/k8s-lab/pkg/cluster"
)

var (
	app         = kingpin.New("k8s-lab", "Workshop/lab tool for managing local Kubernetes")
	clusterName = app.Flag("cluster-name", "Cluster name").Default("k8s-lab").String()

	initCmd = app.Command("init", "Initialize lab environment")

	deleteCmd = app.Command("delete", "Delete lab environment")
)

func main() {
	ctx := context.Background()

	switch cmd := kingpin.MustParse(app.Parse(os.Args[1:])); cmd {
	case initCmd.FullCommand():
		runInit(ctx)
	case deleteCmd.FullCommand():
		runDestroy(ctx)
	}
}

func runInit(_ context.Context) {
	if err := cluster.Init(*clusterName); err != nil {
		logrus.WithError(err).Error("init k8s-lab")
	}
}

func runDestroy(_ context.Context) {
	if err := cluster.Delete(*clusterName); err != nil {
		logrus.WithError(err).Error("delete k8s-lab")
	}
}
