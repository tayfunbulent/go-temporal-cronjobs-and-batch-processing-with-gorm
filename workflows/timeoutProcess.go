package workflows

import (
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func SimulateTimeout(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 3,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	if err := workflow.ExecuteActivity(ctx, activities.SimulateTimeout).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
