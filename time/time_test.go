package time

import (
	"context"

	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	timer := time.NewTicker(time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case now := <-timer.C:
			t.Log(now.String())
		case <-ctx.Done():
			//退出
			return
		}

	}
}

func TestTimer(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	go func() {
		for now := range timer.C {
			t.Log(now.Unix())
		}

	}()
	time.Sleep(10 * time.Second)
}
