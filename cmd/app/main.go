package main

import (
	"go.temporal.io/sdk/client"
	"log"
	temporalClient "temporalpoc/client"
	"temporalpoc/types"
	"temporalpoc/workflows"
)

func main() {
	temporalClient.Init()

	options := client.StartWorkflowOptions{
		ID:        "test",
		TaskQueue: temporalClient.QUEUE_NAME,
	}

	input := types.PaymentInput{
		SrcAccount:    "10",
		TargetAccount: "20",
		Amount:        250,
	}
	err := temporalClient.StartWorkflow(input, workflows.MoneyTransfer, options)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}
}
