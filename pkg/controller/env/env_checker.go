package env

import (
	"context"

	"github.com/opensergo/opensergo-control-plane/pkg/client"
)

const (
	K8S_ENV   = "K8S"
	ISTIO_ENV = "ISTIO"

	ISTIO_DEPLOYMENT_NAME = "istiod"
	ISTIO_NAMESPACE       = "istio-system"
)

func GetENV() string {
	_, err := client.GetDeployment(context.Background(), ISTIO_DEPLOYMENT_NAME, ISTIO_NAMESPACE)
	if err != nil {
		return K8S_ENV
	}
	return ISTIO_ENV
}
