package payment

import (
	"context"
	"log"
	"temporalpoc/types"
)

func Withdraw(ctx context.Context, data types.PaymentInput) (string, error) {
	log.Printf("Withdrawing $%d from account %s.\n\n",
		data.Amount,
		data.SrcAccount,
	)
	return "confirmed", nil
}

func Deposit(ctx context.Context, data types.PaymentInput) (string, error) {
	log.Printf("Depositing $%d into account %s.\n\n",
		data.Amount,
		data.TargetAccount,
	)

	return "", nil
}
