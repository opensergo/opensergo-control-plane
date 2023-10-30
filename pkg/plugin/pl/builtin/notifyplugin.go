package builtin

import (
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"
	ratelimit_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/ratelimit"
)

func NotifyPluginRateLimit(r ratelimit_plugin.RateLimit, l *v1alpha1.RateLimitStrategy) error {
	limit, err := r.RateLimit(l.Spec.Threshold)
	if err != nil {
		return err
	}
	l.Spec.Threshold = limit
	return nil
}
