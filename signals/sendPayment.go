package signals

import (
	"context"
	"go.temporal.io/sdk/client"
	"log"
)

const (
	PaymentSignal = "payment_signal"
)

func SendPayment(workflowID string, paymentStatus bool) (err error) {
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Println("Unable to create Temporal server", err)
		return
	}

	err = temporalClient.SignalWorkflow(context.Background(), workflowID, "", PaymentSignal, paymentStatus)
	if err != nil {
		log.Println("Error signaling server", err)
		return
	}

	return nil
}
