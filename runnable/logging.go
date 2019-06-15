package runnable

import (
	"context"

	"code.cloudfoundry.org/lager"
)

type WithLogging struct {
	runnable Runnable
	logger   lager.Logger
}

var _ Runnable = &WithLogging{}

func NewWithLogging(logger lager.Logger, runnable Runnable) *WithLogging {
	return &WithLogging{
		runnable: runnable,
		logger:   logger,
	}
}

func (r *WithLogging) Run(ctx context.Context) (err error) {
	r.logger.Info("start")

	err = r.runnable.Run(ctx)
	if err != nil {
		r.logger.Error("finish", err)
		return
	}

	r.logger.Info("finish")
	return
}
