package main

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/activities"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Worker Starting...")
	opt := client.Options{
		HostPort: client.DefaultHostPort,
	}
	c, err := client.NewClient(opt)
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Worker 1 - Task Queue: "worker-group-1"
	w1 := worker.New(c, "worker-group-1", worker.Options{})
	w1.RegisterWorkflow(workflows.CoffeeShopWorkflow)
	w1.RegisterActivity(activities.PrepareCoffee)
	w1.RegisterActivity(activities.WriteAsDept)
	w1.RegisterActivity(activities.GiveCoffee)

	// Worker 2 - Task Queue: "worker-group-2"
	w2 := worker.New(c, "worker-group-2", worker.Options{})
	w2.RegisterWorkflow(workflows.TimeoutWorkflow)
	w2.RegisterActivity(activities.LongRunningActivity)

	// Worker 3 - Task Queue: "worker-group-3"
	w3 := worker.New(c, "worker-group-3", worker.Options{})
	w3.RegisterWorkflow(workflows.EmailBatchWorkflow)
	w3.RegisterActivity(activities.SendEmail)

	// Worker'ları ayrı goroutine'lerde çalıştır
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := w1.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Worker 1 stopped with error", err)
		}
	}()

	go func() {
		if err := w2.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Worker 2 stopped with error", err)
		}
	}()

	go func() {
		if err := w3.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Worker 3 stopped with error", err)
		}
	}()

	// JOB
	w := worker.New(c, "worker-group", worker.Options{})
	w.RegisterWorkflow(workflows.EmailJobWorkflow)
	w.RegisterActivity(activities.SendCurlRequest)

	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	}()

	// Programı sinyal alınana kadar beklet
	<-ctx.Done()

	log.Println("Shutting down workers...")
	stop()
}
