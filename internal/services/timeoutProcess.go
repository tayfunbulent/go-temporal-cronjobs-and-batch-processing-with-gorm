package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/starters"
	"net/http"
)

func SimulateTimeout(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SimulateTimeout(workflowID.String())
	return c.String(
		http.StatusOK,
		fmt.Sprintf("%s check to CLI or Temporal Web UI for the status of the workflow", workflowID.String()),
	)
}
