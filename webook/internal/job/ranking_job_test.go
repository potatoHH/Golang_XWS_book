package job

import (
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/pkg/logger"
	"testing"
	"time"
)

func TestRankingJobAdapter_Run(t *testing.T) {
	type fields struct {
		j Job
		l logger.LoggerV1
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RankingJobAdapter{
				j: tt.fields.j,
				l: tt.fields.l,
			}
			r.Run()
		})
	}
}

func TestRankingJob_Name(t *testing.T) {
	type fields struct {
		svc     service.RankingService
		timeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RankingJob{
				svc:     tt.fields.svc,
				timeout: tt.fields.timeout,
			}
			if got := r.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRankingJob_Run(t *testing.T) {
	type fields struct {
		svc     service.RankingService
		timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RankingJob{
				svc:     tt.fields.svc,
				timeout: tt.fields.timeout,
			}
			if err := r.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
