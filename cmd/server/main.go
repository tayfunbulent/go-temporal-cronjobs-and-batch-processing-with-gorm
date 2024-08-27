package main

import (
	"github.com/labstack/echo/v4"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/internal/cronjobs"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/internal/services"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal server", err)
	}
	defer c.Close()

	cronjobs.ConfigureCronJobs(c)

	e := echo.New()
	e.GET("/health-check", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	startWorkflow := e.Group("/start-workflow")
	{
		sandwichShop := startWorkflow.Group("/sandwich-shop")
		{
			sandwichShop.GET("/place-order", services.PlaceOrder)
			sandwichShop.GET("/check-payment/:workflowID", services.CheckPayment)
		}

		sendEmail := startWorkflow.Group("/send-email")
		{
			sendEmail.GET("/v1", services.SimulateSendEmailWithLinearProcess)
			sendEmail.GET("/v2", services.SimulateSendEmailWithTemporalGoroutineAndBatchProcess)
			sendEmail.GET("/v3", services.SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess)
			sendEmail.GET("/v4", services.SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess)
			sendEmail.GET("/all-versions", func(c echo.Context) error {
				_ = c.String(200, "v1 : ")
				if err = services.SimulateSendEmailWithLinearProcess(c); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				_ = c.String(200, "\nv2 : ")
				if err = services.SimulateSendEmailWithTemporalGoroutineAndBatchProcess(c); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				_ = c.String(200, "\nv3 : ")
				if err = services.SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess(c); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				_ = c.String(200, "\nv4 : ")
				if err = services.SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess(c); err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				_ = c.String(200, "\nResult : ")
				return c.String(http.StatusOK, "All workflows started")
			})
		}

		startWorkflow.GET("/timeout", services.SimulateTimeout)
	}

	e.Logger.Fatal(e.Start(":3310"))
}
