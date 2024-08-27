package services

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/starters"
	"net/http"
	"strconv"
)

func SimulateSendEmailWithLinearProcess(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SimulateSendEmailWithLinearProcess(workflowID.String(), generateRandomEmails())
	return c.String(http.StatusOK, workflowID.String())
}

func SimulateSendEmailWithTemporalGoroutineAndBatchProcess(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SimulateSendEmailWithTemporalGoroutineAndBatchProcess(workflowID.String(), generateRandomEmails())
	return c.String(http.StatusOK, workflowID.String())
}

func SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess(workflowID.String(), generateRandomEmails())
	return c.String(http.StatusOK, workflowID.String())
}

func SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess(workflowID.String(), generateRandomEmails())
	return c.String(http.StatusOK, workflowID.String())
}

func generateRandomEmails() []string {
	var emails = make([]string, 100)

	for i := 0; i < 100; i++ {
		emails[i] = strconv.Itoa(i) + "@example.com"
	}

	return emails
}
