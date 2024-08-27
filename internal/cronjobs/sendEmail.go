package cronjobs

import (
	"context"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func sendEmail(c client.Client, frequency string) {
	workflowID := "send_email_job"
	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    "send-email-from-cron-job",
		CronSchedule: frequency,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SimulateSendEmailWithCronJob)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	log.Printf("Started workflow with ID: %s\n", we.GetID())

	resp, err := c.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if err != nil {
		log.Fatalln("Unable to describe workflow execution", err)
	}

	status := resp.WorkflowExecutionInfo.Status
	if status.String() != "Completed" {
		log.Printf("Workflow '%s' is not completed yet. Current status: %s\n", workflowID, status.String())
	} else {
		log.Printf("Workflow '%s' is completed\n", workflowID)
	}

	startTime := resp.WorkflowExecutionInfo.GetStartTime()
	log.Printf("Workflow '%s' last started at: %s\n", workflowID, startTime.AsTime().Format(time.RFC3339))
}
