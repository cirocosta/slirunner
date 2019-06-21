package runnable

import (
	"context"

	"code.cloudfoundry.org/lager/lagerctx"
)

type WithLogging struct {
	runnable Runnable
	name     string
}

var _ Runnable = &WithLogging{}

func NewWithLogging(name string, runnable Runnable) *WithLogging {
	return &WithLogging{
		runnable: runnable,
		name:     name,
	}
}

func (r *WithLogging) Run(ctx context.Context) (err error) {
	logger := lagerctx.FromContext(ctx).Session(r.name)
	logger.Info("start")

	err = r.runnable.Run(ctx)
	if err != nil {
		logger.Error("finish", err)
		return
	}

	logger.Info("finish")
	return
}
