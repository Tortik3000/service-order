package rate_limit

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RateLimitTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rate_limit_total",
		Help: "Total number of rate limit",
	},
	[]string{"method", "url"},
)
