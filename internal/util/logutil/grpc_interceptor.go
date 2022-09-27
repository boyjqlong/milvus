package logutil

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus/internal/util/typeutil"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/milvus-io/milvus/internal/log"
	"github.com/milvus-io/milvus/internal/util/trace"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	logLevelRPCMetaKey = "log_level"
	clientRequestIDKey = "client_request_id"
)

type ctxInjector func(context.Context) context.Context

// UnaryTraceLoggerInterceptor adds a traced logger in unary rpc call ctx
func UnaryTraceLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	opt := addInjectors(withLevelAndTrace)
	// TODO: inject module name for different components.
	return unaryTraceLoggerInterceptor(ctx, req, info, handler, opt)
}

// StreamTraceLoggerInterceptor add a traced logger in stream rpc call ctx
func StreamTraceLoggerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	opt := addInjectors(withLevelAndTrace)
	// TODO: inject module name for different components.
	return streamTraceLoggerInterceptor(srv, ss, info, handler, opt)
}

func ProxyUnaryTraceLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	opt := proxyConfigOpt(ctx)
	return unaryTraceLoggerInterceptor(ctx, req, info, handler, opt)
}

func unaryTraceLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, opts ...configOpt) (interface{}, error) {
	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}

	newCtx := c.apply(ctx)

	c.runBeforeHandler()
	resp, err := handler(newCtx, req)
	c.runAfterHandler()

	return resp, err
}

func streamTraceLoggerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler, opts ...configOpt) error {
	ctx := ss.Context()

	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}

	newCtx := c.apply(ctx)

	wrappedStream := grpc_middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = newCtx

	c.runBeforeHandler()
	err := handler(srv, wrappedStream)
	c.runAfterHandler()

	return err
}

type config struct {
	before       []func()
	after        []func()
	ctxInjectors []ctxInjector
}

func (c config) apply(ctx context.Context) context.Context {
	newCtx := ctx
	for _, injector := range c.ctxInjectors {
		newCtx = injector(newCtx)
	}
	return newCtx
}

func (c config) runBeforeHandler() {
	for _, cb := range c.before {
		cb()
	}
}

func (c config) runAfterHandler() {
	for _, cb := range c.after {
		cb()
	}
}

func (c *config) addBeforeCallbacks(callbacks ...func()) {
	c.before = append(c.before, callbacks...)
}

func (c *config) addAfterCallbacks(callbacks ...func()) {
	c.after = append(c.after, callbacks...)
}

func (c *config) addInjectors(injectors ...ctxInjector) {
	c.ctxInjectors = append(c.ctxInjectors, injectors...)
}

func defaultConfig() *config {
	return &config{
		before:       make([]func(), 0),
		after:        make([]func(), 0),
		ctxInjectors: make([]ctxInjector, 0),
	}
}

type configOpt func(c *config)

func addBeforeCallbacks(callbacks ...func()) configOpt {
	return func(c *config) {
		c.addBeforeCallbacks(callbacks...)
	}
}

func addAfterCallbacks(callbacks ...func()) configOpt {
	return func(c *config) {
		c.addAfterCallbacks(callbacks...)
	}
}

func addInjectors(injectors ...ctxInjector) configOpt {
	return func(c *config) {
		c.addInjectors(injectors...)
	}
}

func withLevelAndTrace(ctx context.Context) context.Context {
	newctx := ctx
	var traceID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		levels := md.Get(logLevelRPCMetaKey)
		// get log level
		if len(levels) >= 1 {
			level := zapcore.DebugLevel
			if err := level.UnmarshalText([]byte(levels[0])); err != nil {
				newctx = ctx
			} else {
				switch level {
				case zapcore.DebugLevel:
					newctx = log.WithDebugLevel(ctx)
				case zapcore.InfoLevel:
					newctx = log.WithInfoLevel(ctx)
				case zapcore.WarnLevel:
					newctx = log.WithWarnLevel(ctx)
				case zapcore.ErrorLevel:
					newctx = log.WithErrorLevel(ctx)
				case zapcore.FatalLevel:
					newctx = log.WithFatalLevel(ctx)
				default:
					newctx = ctx
				}
			}
			// inject log level to outgoing meta
			newctx = metadata.AppendToOutgoingContext(newctx, logLevelRPCMetaKey, level.String())
		}
		// client request id
		requestID := md.Get(clientRequestIDKey)
		if len(requestID) >= 1 {
			traceID = requestID[0]
			// inject traceid in order to pass client request id
			newctx = metadata.AppendToOutgoingContext(newctx, clientRequestIDKey, traceID)
		}
	}
	if traceID == "" {
		traceID, _, _ = trace.InfoFromContext(newctx)
	}
	if traceID != "" {
		newctx = log.WithTraceID(newctx, traceID)
	}
	return newctx
}

func withLevel(ctx context.Context) context.Context {
	newctx := ctx
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		levels := md.Get(logLevelRPCMetaKey)
		// get log level
		if len(levels) >= 1 {
			level := zapcore.DebugLevel
			if err := level.UnmarshalText([]byte(levels[0])); err != nil {
				newctx = ctx
			} else {
				switch level {
				case zapcore.DebugLevel:
					newctx = log.WithDebugLevel(ctx)
				case zapcore.InfoLevel:
					newctx = log.WithInfoLevel(ctx)
				case zapcore.WarnLevel:
					newctx = log.WithWarnLevel(ctx)
				case zapcore.ErrorLevel:
					newctx = log.WithErrorLevel(ctx)
				case zapcore.FatalLevel:
					newctx = log.WithFatalLevel(ctx)
				default:
					newctx = ctx
				}
			}
			// inject log level to outgoing meta
			newctx = metadata.AppendToOutgoingContext(newctx, logLevelRPCMetaKey, level.String())
		}
	}
	return newctx
}

func withProxyModule(ctx context.Context) context.Context {
	return log.WithModule(ctx, typeutil.ProxyRole)
}

func injectTraceID(ctx context.Context, traceID string) context.Context {
	newCtx := metadata.AppendToOutgoingContext(ctx, clientRequestIDKey, traceID)
	return newCtx
}

func getClientTraceID(ctx context.Context) (has bool, traceID string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// client request id
		requestID := md.Get(clientRequestIDKey)
		if len(requestID) >= 1 {
			traceID = requestID[0]
			// inject traceid in order to pass client request id
			return true, traceID
		}
	}
	return false, ""
}

func proxyConfigOpt(ctx context.Context) configOpt {
	hasTraceID, traceID := getClientTraceID(ctx)
	return func(c *config) {
		span, newCtx := trace.StartSpanFromContext(ctx)
		c.addInjectors(func(ctx context.Context) context.Context {
			return newCtx
		})
		c.addInjectors(withLevel)
		c.addInjectors(withProxyModule)
		if hasTraceID {
			c.addInjectors(func(ctx context.Context) context.Context {
				return injectTraceID(ctx, traceID)
			})
		} else {
			traceID, _, _ = trace.InfoFromSpan(span)
		}
		c.addInjectors(func(ctx context.Context) context.Context {
			return log.WithTraceID(ctx, traceID)
		})
		c.addBeforeCallbacks(func() {
			fmt.Println("test if before hook worked: ", traceID)
		})
		c.addAfterCallbacks(
			func() { span.Finish() },
			func() { fmt.Println("test if after hook worked: ", traceID) },
		)
	}
}
