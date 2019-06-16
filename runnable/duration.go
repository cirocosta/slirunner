package runnable

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type WithMetrics struct {
	runnable Runnable
	name     string
}

var (
	_ Runnable = &WithMetrics{}

	runs = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "slirunner_probe_duration_seconds",
		Help:    "Time spent running a specific probe",
		Buckets: []float64{1, 5, 10, 20, 30, 40, 50, 60, 90, 120},
	}, []string{"probe", "result"})
)

func NewWithMetrics(probeName string, runnable Runnable) *WithMetrics {
	return &WithMetrics{
		runnable: runnable,
		name:     probeName,
	}
}

func (r *WithMetrics) Run(ctx context.Context) (err error) {
	start := time.Now()

	defer func() {
		result := "success"
		if err != nil {
			result = "failure"
		}

		runs.WithLabelValues(r.name, result).Observe(time.Since(start).Seconds())
	}()

	err = r.runnable.Run(ctx)
	if err != nil {
		return
	}

	return
}
