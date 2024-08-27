package signals

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

func ReceiveSignal(ctx workflow.Context, signalName string, timeout time.Duration) (paymentStatus bool) {
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	timer := workflow.NewTimer(ctx, timeout)
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &paymentStatus)
	}).AddFuture(timer, func(f workflow.Future) {
		paymentStatus = false
	})

	selector.Select(ctx)

	return
}
