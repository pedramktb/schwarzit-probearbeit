package usecase

import "context"

type PaymentStatusChangeUC interface {
	PaymentStatusChange(ctx context.Context, id string) error
}
