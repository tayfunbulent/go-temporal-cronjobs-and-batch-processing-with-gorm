package starters

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
)

func SimulateSendEmailWithLinearProcess(workflowID string, emails []string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "send-email-with-linear-process",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SimulateSendEmailWithLinearProcess, emails); err != nil {
		panic(err)
	}
}

func SimulateSendEmailWithTemporalGoroutineAndBatchProcess(workflowID string, emails []string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "send-email-with-goroutine-and-batch-process",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SimulateSendEmailWithTemporalGoroutineAndBatchProcess, emails); err != nil {
		panic(err)
	}
}

func SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess(workflowID string, emails []string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "send-email-with-child-and-batch-process",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess, emails); err != nil {
		panic(err)
	}
}

func SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess(workflowID string, emails []string) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	opt := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "send-email-with-child-and-future-and-batch-process",
	}
	if _, err := c.ExecuteWorkflow(context.Background(), opt, workflows.SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess, emails); err != nil {
		panic(err)
	}
}
