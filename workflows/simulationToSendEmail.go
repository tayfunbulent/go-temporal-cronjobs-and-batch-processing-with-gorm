package workflows

import (
	"fmt"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"math"
	"runtime"
	"time"
)

func measureMemoryUsage() uint64 {
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc
}

type EmailResult struct {
	Counter int
}

func SimulateSendEmailWithLinearProcess(ctx workflow.Context, emailAddresses []string) error {
	startTime := time.Now()
	initialMemory := measureMemoryUsage()

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

	counter := 0

	for _, email := range emailAddresses {
		if err := workflow.ExecuteActivity(ctx, activities.SendEmail, email).Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Failed to send email", "email", email, "error", err)
		} else {
			counter++
		}
	}

	duration := time.Since(startTime)
	finalMemory := measureMemoryUsage()
	fmt.Printf("SimulateSendEmailWithLinearProcess | %d emails sent and batches processed in %s\n", counter, duration)
	fmt.Printf("SimulateSendEmailWithLinearProcess | Memory used: %d bytes\n", finalMemory-initialMemory)

	return nil
}

func SimulateSendEmailWithTemporalGoroutineAndBatchProcess(ctx workflow.Context, emailAddresses []string) error {
	startTime := time.Now()
	initialMemory := measureMemoryUsage()

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
				if err := workflow.ExecuteActivity(ctx, activities.SendEmail, email).Get(ctx, nil); err != nil {
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

	duration := time.Since(startTime)
	finalMemory := measureMemoryUsage()
	fmt.Printf("SimulateSendEmailWithTemporalGoroutineAndBatchProcess | %d emails sent and batches processed in %s\n", counter, duration)
	fmt.Printf("SimulateSendEmailWithTemporalGoroutineAndBatchProcess | Memory used: %d bytes\n", finalMemory-initialMemory)

	return nil
}

func SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess(ctx workflow.Context, emailAddresses []string) error {
	startTime := time.Now()
	initialMemory := measureMemoryUsage()

	batchSize := 10
	totalBatches := int(math.Ceil(float64(len(emailAddresses)) / float64(batchSize)))

	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: time.Hour,
	}

	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	totalCounter := 0
	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize

		if end > len(emailAddresses) {
			end = len(emailAddresses)
		}

		batch := emailAddresses[start:end]
		var result *EmailResult
		if err := workflow.ExecuteChildWorkflow(ctx, SimulateSendEmailChildWorkflow, batch).Get(ctx, &result); err != nil {
			return err
		} else {
			totalCounter += result.Counter
		}
	}

	duration := time.Since(startTime)
	finalMemory := measureMemoryUsage()
	fmt.Printf("SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess | %d emails sent and batches processed in %s\n", totalCounter, duration)
	fmt.Printf("SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess | Memory used: %d bytes\n", finalMemory-initialMemory)

	return nil
}

func SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess(ctx workflow.Context, emailAddresses []string) error {
	startTime := time.Now()
	initialMemory := measureMemoryUsage()

	batchSize := 10
	totalBatches := int(math.Ceil(float64(len(emailAddresses)) / float64(batchSize)))

	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: time.Hour,
	}

	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	var futures []workflow.Future

	for i := 0; i < totalBatches; i++ {
		start := i * batchSize
		end := start + batchSize

		if end > len(emailAddresses) {
			end = len(emailAddresses)
		}

		batch := emailAddresses[start:end]

		future := workflow.ExecuteChildWorkflow(ctx, SimulateSendEmailChildWorkflow, batch)
		futures = append(futures, future)
	}

	totalCounter := 0
	var result *EmailResult
	for _, future := range futures {
		err := future.Get(ctx, &result)
		if err != nil {
			return err
		} else {
			totalCounter += result.Counter
		}
	}

	duration := time.Since(startTime)
	finalMemory := measureMemoryUsage()
	fmt.Printf("SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess | %d emails sent and batches processed in %s\n", totalCounter, duration)
	fmt.Printf("SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess | Memory used: %d bytes\n", finalMemory-initialMemory)

	return nil
}

func SimulateSendEmailChildWorkflow(ctx workflow.Context, batch []string) (*EmailResult, error) {
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
	counter := 0

	for _, email := range batch {
		if err := workflow.ExecuteActivity(ctx, activities.SendEmail, email).Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Failed to send email", "email", email, "error", err)
			return nil, err
		} else {
			counter++
		}
	}

	return &EmailResult{Counter: counter}, nil
}

func SimulateSendEmailWithCronJob(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	/* different way to send requests with a sleep for 2 minutes | CRONJOB
	for {
		err := workflow.ExecuteActivity(ctx, activities.SendRequestToURL, url).Get(ctx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("Failed to send request", "error", err)
		}

		workflow.Sleep(ctx, 2*time.Minute)
	} */

	endpoint := "http://localhost:3310/start-workflow/send-email/all-versions"

	if err := workflow.ExecuteActivity(ctx, activities.SendRequestToURL, endpoint).Get(ctx, nil); err != nil {
		workflow.GetLogger(ctx).Error("Failed to send request", "error", err)
		return err
	}

	return nil
}
