package memory

import (
	"context"
	"fmt"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//模拟发短信的过程
	fmt.Println(args)
	return nil
}
