package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/signals"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/starters"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/workflows"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CronSchedule struct {
	Schedule string
}

var (
	EveryMinute   = CronSchedule{"* * * * *"}   // Her dakika
	Every2Minutes = CronSchedule{"*/2 * * * *"} // Her 2 dakikada bir
	EveryHour     = CronSchedule{"0 * * * *"}   // Her saat başı
	EveryDay      = CronSchedule{"0 0 * * *"}   // Her gün saat 00:00'da
	EveryWeek     = CronSchedule{"0 0 * * 0"}   // Her hafta Pazar günü saat 00:00'da
	EveryMonth    = CronSchedule{"0 0 1 * *"}   // Her ayın 1'inde saat 00:00'da
)

func main() {
	log.Println("Client Starting...")

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:           "email_job_workflow_v2",
		TaskQueue:    "worker-group",
		CronSchedule: Every2Minutes.Schedule,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.EmailJobWorkflow)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	log.Printf("Started workflow with ID: %s\n", we.GetID())

	workflowID := "email_job_workflow_v2"

	// İş akışının durumunu al
	resp, err := c.DescribeWorkflowExecution(context.Background(), workflowID, "")
	if err != nil {
		log.Fatalln("Unable to describe workflow execution", err)
	}

	// İş akışının durumu
	// status := resp.WorkflowExecutionInfo.Status

	// İş akışının tamamlanıp tamamlanmadığını kontrol et
	/* if status.String() != "Completed" {
		log.Printf("Workflow '%s' is not completed yet. Current status: %s\n", workflowID, status.String())
	} else { */
	// Son başlatılma zamanını al
	startTime := resp.WorkflowExecutionInfo.GetStartTime()

	// Zamanı formatlayarak yazdır
	log.Printf("Workflow '%s' last started at: %s\n", workflowID, time.Unix(0, startTime.UnixNano()).Format(time.RFC3339))
	//}

	e := echo.New()
	e.GET("/about", GetAbout)
	e.GET("/start-workflow-timeout", StartWorkflowFuncTimeout)
	e.GET("/start-workflow-email", StartWorkflowFuncEmail)
	e.GET("/start-workflow", StartWorkflow)
	e.GET("/send-signal/:workflowID", SendSignal)

	e.Logger.Fatal(e.Start(":3310"))

}

func StartWorkflow(c echo.Context) error {
	workflowUUID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.StartWorkflowFunc(workflowUUID.String(), "test-customer-name")
	return c.String(http.StatusOK, workflowUUID.String())
}

func StartWorkflowFuncTimeout(c echo.Context) error {
	workflowUUID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.StartWorkflowFuncTimeout(workflowUUID.String())
	return c.String(http.StatusOK, workflowUUID.String())
}

func StartWorkflowFuncEmail(c echo.Context) error {
	workflowUUID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var emails = make([]string, 100)

	for i := 0; i < 100; i++ {
		emails[i] = strconv.Itoa(i) + "@example.com"
	}

	starters.StartWorkflowFuncEmail(workflowUUID.String(), emails)
	return c.String(http.StatusOK, workflowUUID.String())
}

func GetAbout(c echo.Context) error {
	return c.String(http.StatusOK, "It's running...")
}

func SendSignal(c echo.Context) error {
	workflowId := c.Param("workflowID")
	if err := signals.SendPaymentSignal(workflowId, false); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Sent it")
}
