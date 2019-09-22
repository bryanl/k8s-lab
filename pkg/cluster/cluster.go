package cluster

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/cluster"
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

	context := cluster.NewContext(clusterName)

	logger.Info("init cluster")

	if err := context.Create(); err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}

	return nil
}

func Delete(clusterName string) error {
	logger := logrus.WithField("cluster-name", clusterName)

	context := cluster.NewContext(clusterName)

	logger.Info("delete cluster")

	if err := context.Delete(); err != nil {
		return fmt.Errorf("delete cluster: %w", err)
	}

	return nil
}
