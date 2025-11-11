package service

import (
	domain2 "Book_Exp/webook/interactive/domain"
	service2 "Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/domain"
	svcmocks "Book_Exp/webook/internal/service/mocks"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRankingTopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name      string
		batchSize int
		n         int
		mock      func(ctrl *gomock.Controller) (ArticleService, service2.InteractiveService)
		wantErr   error
		wantArt   []domain.Article
	}{
		{
			name: "计算成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, service2.InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				artSvc.EXPECT().ListPub(gomock.Any(), 0, 3).Return([]domain.Article{
					{Id: 1, Utime: now, Ctime: now},
					{Id: 2, Utime: now, Ctime: now},
					{Id: 3, Utime: now, Ctime: now},
				}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), 3, 3).
					Return([]domain.Article{}, nil)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					"article", []int64{1, 2, 3}).
					Return(map[int64]domain2.Interactive{
						1: {BizId: 1, LikeCnt: 1},
						2: {BizId: 2, LikeCnt: 2},
						3: {BizId: 3, LikeCnt: 3},
					}, nil)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					"article", []int64{}).
					Return(map[int64]domain2.Interactive{}, nil)
				return artSvc, intrSvc
			},
			wantArt: []domain.Article{
				{Id: 1, Utime: now, Ctime: now},
				{Id: 2, Utime: now, Ctime: now},
				{Id: 3, Utime: now, Ctime: now},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, intrSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(artSvc, intrSvc)

			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantArt, arts)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
