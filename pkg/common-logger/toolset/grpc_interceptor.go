package toolset

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

func GRPCClientTracingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if val := ctx.Value(logger.TraceIDKey).(string); val == "" {
		ctx = metadata.AppendToOutgoingContext(ctx, string(logger.TraceIDKey), logger.GetNewTraceID())
	}

	// Calls the invoker to execute RPC
	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}

func GRPCServerTracingInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		traceID   string
		tracedCtx context.Context
	)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		trace := md.Get(string(logger.TraceIDKey))
		if len(trace) != 0 {
			traceID = trace[0]
		} else {
			traceID = logger.GetNewTraceID()
		}
	} else {
		traceID = logger.GetNewTraceID()
	}

	tracedCtx = context.WithValue(ctx, logger.TraceIDKey, traceID)

	return handler(tracedCtx, req)
}
