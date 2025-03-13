package client

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vucongthanh92/go-base-utils/grpc/interceptors"
	"github.com/vucongthanh92/go-base-utils/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewClientConn(ctx context.Context, logger logger.Logger, port string, development bool) *grpc.ClientConn {
	unaryInterceptorOption := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
				grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
				grpc_retry.WithMax(backoffRetries),
			),
			otelgrpc.UnaryClientInterceptor(),
			interceptors.ClientLogger(logger),
		))

	streamInterceptors := grpc.WithStreamInterceptor(
		otelgrpc.StreamClientInterceptor(),
	)

	var opts grpc.DialOption
	if development {
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	clientConn, err := grpc.DialContext(
		ctx,
		port,
		unaryInterceptorOption,
		streamInterceptors,
		opts,
	)

	if err != nil {
		logger.Error("NewClientConn", zap.Error(err))
		return nil
	}

	return clientConn
}
