package signals

import (
	"context"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"log"
	"time"
)

const (
	PAYMENT_SIGNAL = "payment_signal"
)

func SendPaymentSignal(workflowID string, paymentStatus bool) (err error) {
	temporalClient, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
		return
	}

	err = temporalClient.SignalWorkflow(context.Background(), workflowID, "", PAYMENT_SIGNAL, paymentStatus)
	if err != nil {
		log.Fatalln("Error signaling client", err)
		return
	}

	return nil
}

func ReciveSignal(ctx workflow.Context, signalName string, timeout time.Duration) (paymentStatus bool) {
	// Sinyal kanalını alın
	signalChan := workflow.GetSignalChannel(ctx, signalName)

	// Zaman aşımı için bir timer oluşturun
	timer := workflow.NewTimer(ctx, timeout)

	// Sinyal veya zaman aşımı arasında seçim yapmak için bir selector kullanın
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &paymentStatus)
	}).AddFuture(timer, func(f workflow.Future) {
		// Eğer zaman aşımı gerçekleşirse, paymentStatus varsayılan olarak false kalır
		paymentStatus = false
	})

	// Sinyal ya da zaman aşımı gelene kadar bekle
	selector.Select(ctx)

	return
}
