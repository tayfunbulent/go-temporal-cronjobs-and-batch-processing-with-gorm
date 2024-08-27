package starters

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
)

func SandwichShop(workflowID string, customerName string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "sandwich-shop",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SandwichShop, customerName); err != nil {
		panic(err)
	}
}
