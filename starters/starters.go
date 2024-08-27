package starters

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
)

func StartWorkflowFunc(workflowID string, customerName string) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "worker-group-1",
	}
	ctx := context.Background()
	if _, err := c.ExecuteWorkflow(ctx, opt, workflows.CoffeeShopWorkflow, customerName); err != nil {
		panic(err)
	}
}

func StartWorkflowFuncTimeout(workflowID string) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "worker-group-2",
	}
	ctx := context.Background()
	if _, err := c.ExecuteWorkflow(ctx, opt, workflows.TimeoutWorkflow); err != nil {
		panic(err)
	}
}

func StartWorkflowFuncEmail(workflowID string, emails []string) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "worker-group-3",
	}
	ctx := context.Background()
	if _, err := c.ExecuteWorkflow(ctx, opt, workflows.EmailBatchWorkflow, emails); err != nil {
		panic(err)
	}
}
