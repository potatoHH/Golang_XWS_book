package logging

import (
	"Book_Exp/webook/pkg/logger"
	"context"
	"fmt"
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	l logger.LoggerV1
	//fn func(msg string, fields ...logger.Field)
}

func (i *InterceptorBuilder) Builder() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		var event = "normal"
		defer func() {
			//持续时间
			duration := time.Since(start)
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				event = "recover"
				err = status.New(codes.Internal, "panic,err"+err.Error()).Err()

			}

			fields := []logger.Field{
				logger.Int64("duration", duration.Milliseconds()),
				logger.String("type", "unary"),
				logger.String("event", event),
				logger.String("method", info.FullMethod),
				//这一部分是你客户端配合的

			}
			if err != nil {
				st, _ := status.FromError(err)
				fields = append(fields, logger.String("code", st.Code().String()),
					logger.String("msg_code", st.Message()))
			}
			i.l.Info("RPC请求", fields...)

		}()
		resp, err = handler(ctx, req)
		return
	}

}
