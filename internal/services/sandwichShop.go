package services

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/signals"
	"github.com/tayfunbulent/go-temporal-cronjobs-and-batch-processing-with-gorm/starters"
	"net/http"
)

func PlaceOrder(c echo.Context) error {
	workflowID, err := uuid.NewUUID()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	starters.SandwichShop(workflowID.String(), "customerName")
	return c.String(http.StatusOK, workflowID.String())
}

func CheckPayment(c echo.Context) error {
	workflowID := c.Param("workflowID")
	if err := signals.SendPayment(workflowID, true); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return c.String(http.StatusNotFound, "No records found for the given workflow ID")
		case "workflow execution already completed":
			return c.String(http.StatusBadRequest, "Workflow execution already completed")
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.String(http.StatusOK, "Payment completed")
}
