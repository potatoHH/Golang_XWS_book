package trace

import (
	"Book_Exp/webook/pkg/grpcx/interceptors"
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//链路控制

type InterceptorBuilder struct {
	propagator propagation.TextMapPropagator
	tracer     trace.Tracer
	interceptors.Builder
}

func (b *InterceptorBuilder) BuildClient() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		//inject
		//要把和trace有关的链路元素据, 传递到服务端
		propagator := b.propagator
		if propagator == nil {
			//这个是全局的
			propagator = otel.GetTextMapPropagator()
		}
		tracer := b.tracer
		if tracer == nil {
			tracer = otel.GetTracerProvider().Tracer("grpcx")
		}
		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String("grpc"),
			attribute.Key("rpc.grpc.kind").String("unary"),
			attribute.Key("rpc.component").String("client"),
		}
		ctx, span := tracer.Start(ctx, method, trace.WithAttributes(attrs...))
		defer span.End()
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
		}()
		ctx = inject(ctx, propagator)
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}

}

func (b *InterceptorBuilder) BuildServer() grpc.UnaryServerInterceptor {
	propagator := b.propagator
	if propagator == nil {
		//这个是全局的
		propagator = otel.GetTextMapPropagator()
	}
	tracer := b.tracer
	if tracer == nil {
		tracer = otel.GetTracerProvider().Tracer("grpcx")
	}
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		attribute.Key("rpc.grpc.kind").String("unary"),
		attribute.Key("rpc.component").String("server"),
	}
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = extract(ctx, propagator)
		ctx, span := b.tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer func() {
			span.End()
		}()
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(b.PeerName(ctx)),
			attribute.Key("net.peer.ip").String(b.PeerIP(ctx)),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := errors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
		}()
		return handler(ctx, req)
	}

}

func inject(ctx context.Context, propagators propagation.TextMapPropagator) context.Context {
	//先看ctx里面有没有原始数据
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	//把元素据放回去 ctx, 具体什么格式由propagator
	propagators.Inject(ctx, GrpcHeaderCarrier(md))
	//再次封装 后续可以在invoker里面继续拿到metadata
	return metadata.NewOutgoingContext(ctx, md)
}

func extract(ctx context.Context, p propagation.TextMapPropagator) context.Context {
	//这里拿到客户端过来的链路数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	//把md注入到ctx中
	//他的注入方式不同
	return p.Extract(ctx, GrpcHeaderCarrier(md))

}

type GrpcHeaderCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc GrpcHeaderCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// Set stores the key-value pair.
func (mc GrpcHeaderCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc GrpcHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}
