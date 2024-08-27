package workflows

import (
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/signals"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func SandwichShop(ctx workflow.Context, customerName string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Minute,
		ScheduleToCloseTimeout: 100 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    10,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	if err := workflow.ExecuteActivity(ctx, activities.MakeSandwich, nil).Get(ctx, nil); err != nil {
		return err
	}
	if err := workflow.Sleep(ctx, 3*time.Second); err != nil {
		return err
	}

	status := signals.ReceiveSignal(ctx, signals.PaymentSignal, 10*time.Second)
	if !status {
		if err := workflow.ExecuteActivity(ctx, activities.BorrowSandwich, customerName).Get(ctx, nil); err != nil {
			return err
		}
	} else {
		if err := workflow.ExecuteActivity(ctx, activities.ServeSandwich, customerName).Get(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}
