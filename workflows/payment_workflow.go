package workflows

import (
	"fmt"
	"temporalpoc/activities/payment"
	"temporalpoc/types"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoneyTransfer(ctx workflow.Context, input types.PaymentInput) (string, error) {

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        2, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Withdraw money.
	var withdrawOutput string
	withdrawErr := workflow.ExecuteActivity(ctx, payment.Withdraw, input).Get(ctx, &withdrawOutput)
	if withdrawErr != nil {
		fmt.Println("withdrawErr", withdrawErr)
		return "", withdrawErr
	}

	// Deposit money.
	var depositOutput string
	depositErr := workflow.ExecuteActivity(ctx, payment.Deposit, input).Get(ctx, &depositOutput)
	if depositErr != nil {
		fmt.Println("withdrawErr", withdrawErr)
		return "", fmt.Errorf("failed to deposit output: %w", depositErr)
	}

	result := fmt.Sprintf("Transfer complete (transaction IDs: %s, %s)", withdrawOutput, depositOutput)
	return result, nil
}
