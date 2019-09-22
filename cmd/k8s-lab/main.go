package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bryanl/k8s-lab/pkg/cluster"
)

var (
	app = kingpin.New("k8s-lab", "Workshop/lab tool for managing local Kubernetes")

	initCmd         = app.Command("init", "Initialize lab environment")
	initClusterName = initCmd.Flag("cluster-name", "ClusterName").Default("k8s-lab").String()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case initCmd.FullCommand():
		runInit()
	}
}

func runInit() {
	if err := cluster.Init(*initClusterName); err != nil {
		logrus.WithError(err).Error("init k8s-lab")
	}
}
