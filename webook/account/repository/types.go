package repository

import (
	"Book_Exp/webook/account/domain"
	"context"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
