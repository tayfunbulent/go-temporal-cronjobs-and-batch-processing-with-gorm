package workflows

import (
	"fmt"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/signals"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"log"
	"math"
	"time"
)

func CoffeeShopWorkflow(ctx workflow.Context, customerName string) error {
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

	log.Println(customerName, "is ordering coffee")

	workflow.ExecuteActivity(ctx, activities.PrepareCoffee, nil).Get(ctx, nil)
	workflow.Sleep(ctx, 5*time.Second)
	log.Println("Coffee is ready to serve")

	if status := signals.ReciveSignal(ctx, signals.PAYMENT_SIGNAL, 3*time.Second); !status {
		log.Println("Coffee couldn't be completed! ")
		workflow.ExecuteActivity(ctx, activities.WriteAsDept, customerName).Get(ctx, nil)
	}

	workflow.ExecuteActivity(ctx, activities.GiveCoffee, customerName).Get(ctx, nil)

	log.Println("Coffee done! ")

	return nil
}

func EmailBatchWorkflow(ctx workflow.Context, emailAddresses []string) error {
	batchSize := 10
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	totalBatches := int(math.Ceil(float64(len(emailAddresses)) / float64(batchSize)))
	doneCh := workflow.NewChannel(ctx)

	counter := 0

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize

		if end > len(emailAddresses) {
			end = len(emailAddresses)
		}

		batch := emailAddresses[start:end]

		workflow.Go(ctx, func(ctx workflow.Context) {
			for _, email := range batch {
				err := workflow.ExecuteActivity(ctx, activities.SendEmail, email).Get(ctx, nil)
				if err != nil {
					workflow.GetLogger(ctx).Error("Failed to send email", "email", email, "error", err)
				} else {
					counter++
				}
			}

			doneCh.Send(ctx, true)
		})
	}

	for i := 0; i < totalBatches; i++ {
		doneCh.Receive(ctx, nil)
	}

	fmt.Println("All emails sent - total: ", counter, " emails sent")

	return nil
}

func TimeoutWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 3,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, activities.LongRunningActivity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func EmailJobWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	/* old version
	for {
		err := workflow.ExecuteActivity(ctx, activities.SendCurlRequest, url).Get(ctx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("Failed to send request", "error", err)
		}
		// 10 dakika bekleme
		workflow.Sleep(ctx, 2*time.Minute)
	} */

	err := workflow.ExecuteActivity(ctx, activities.SendCurlRequest).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Failed to send request", "error", err)
		return err
	}

	return nil
}
