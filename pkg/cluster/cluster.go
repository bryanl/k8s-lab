package cluster

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/kind/pkg/cluster"
)

func Init(clusterName string) error {
	isKnown, err := cluster.IsKnown(clusterName)
	if err != nil {
		return fmt.Errorf("check if cluster %s is known: %w", clusterName, err)
	}

	if isKnown {
		logrus.WithField("cluster-name", clusterName).Info("cluster has been previously initialized")
		return nil
	}

	context := cluster.NewContext(clusterName)

	if err := context.Create(); err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}

	return nil
}
