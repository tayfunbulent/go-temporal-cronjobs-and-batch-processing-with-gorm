package main

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/internal/zapadapter"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// WorkerDefinition holds information about each worker
type WorkerDefinition struct {
	TaskQueue         string
	RegisterFunctions func(w worker.Worker)
}

func main() {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)

	logger, err := config.Build()
	if err != nil {
		log.Fatalln("Unable to create logger", err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   zapadapter.NewZapAdapter(logger),
	})

	if err != nil {
		log.Fatalln("Unable to create Temporal server", err)
	}
	defer c.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, wd := range configureWorkers() {
		go func(wd WorkerDefinition) {
			w := worker.New(c, wd.TaskQueue, worker.Options{})
			wd.RegisterFunctions(w)
			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Fatalln("Worker stopped with error", err)
			}
		}(wd)
	}

	<-ctx.Done()

	stop()
}

func configureWorkers() []WorkerDefinition {
	return []WorkerDefinition{
		{
			TaskQueue: "send-email-from-cron-job",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateSendEmailWithCronJob)
				w.RegisterActivity(activities.SendRequestToURL)
			},
		},
		{
			TaskQueue: "sandwich-shop",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SandwichShop)
				w.RegisterActivity(activities.MakeSandwich)
				w.RegisterActivity(activities.BorrowSandwich)
				w.RegisterActivity(activities.ServeSandwich)
			},
		},
		{
			TaskQueue: "send-email-with-linear-process",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateSendEmailWithLinearProcess)
				w.RegisterActivity(activities.SendEmail)
			},
		},
		{
			TaskQueue: "send-email-with-goroutine-and-batch-process",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateSendEmailWithTemporalGoroutineAndBatchProcess)
				w.RegisterActivity(activities.SendEmail)
			},
		},
		{
			TaskQueue: "send-email-with-child-and-batch-process",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess)
				w.RegisterWorkflow(workflows.SimulateSendEmailChildWorkflow)
				w.RegisterActivity(activities.SendEmail)
			},
		},
		{
			TaskQueue: "send-email-with-child-and-future-and-batch-process",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess)
				w.RegisterWorkflow(workflows.SimulateSendEmailChildWorkflow)
				w.RegisterActivity(activities.SendEmail)
			},
		},
		{
			TaskQueue: "simulate-timeout",
			RegisterFunctions: func(w worker.Worker) {
				w.RegisterWorkflow(workflows.SimulateTimeout)
				w.RegisterActivity(activities.SimulateTimeout)
			},
		},
	}
}
