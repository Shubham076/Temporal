package main

import (
	"go.temporal.io/sdk/worker"
	"sync"
	"temporalpoc/activities/payment"
	temporalClient "temporalpoc/client"
	"temporalpoc/workflows"
)

func main() {
	temporalClient.Init()
	wg := &sync.WaitGroup{}
	options := worker.Options{}
	temporalClient.StartWorkers(wg, 2, temporalClient.QUEUE_NAME, options, []interface{}{workflows.MoneyTransfer}, []interface{}{payment.Withdraw, payment.Deposit})
	wg.Wait()
}
