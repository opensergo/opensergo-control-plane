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

var env = ""

func GetENV() string {
	if env != "" {
		return env
	}
	_, err := client.GetDeployment(context.Background(), ISTIO_DEPLOYMENT_NAME, ISTIO_NAMESPACE)
	if err != nil {
		env = K8S_ENV
		return env
	}
	env = ISTIO_ENV
	return env
}
