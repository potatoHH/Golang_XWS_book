package service

import (
	"Book_Exp/webook/account/domain"
	"context"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
