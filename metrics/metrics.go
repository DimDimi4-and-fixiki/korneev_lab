package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MetricVirusesCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "num_viruses_metric",
		Help: "number of viruses found",
	})
)
