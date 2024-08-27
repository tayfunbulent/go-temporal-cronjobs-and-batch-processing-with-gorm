package starters

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
)

func SimulateTimeout(workflowID string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "simulate-timeout",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SimulateTimeout); err != nil {
		panic(err)
	}
}
