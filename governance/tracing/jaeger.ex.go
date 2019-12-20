package tracing

import (
	"context"
	"github.com/joyous-x/saturn/common/xlog"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	tracingCallerAddr = "tracing:calleraddr"
	tracingTargetAddr = "tracing:targetaddr"
)

type textMapCarrier struct {
	metadata.MD
}

func (c textMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c textMapCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}

type jaegerGrpcFilter func(fullMethod string) bool

func jaegerGrpcFilterFunc(fullMethod string) bool {
	if strings.ToLower(fullMethod) == "/grpc.health.v1.health/check" {
		return true
	}
	return false
}

// JaegerGlobalTracer ...
func JaegerGlobalTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

// NewJaegerTracerEx ...
func NewJaegerTracerEx(svc, agentAddr, sampleType, sampleParm string) (opentracing.Tracer, io.Closer) {
	//> 可以借助 "github.com/uber/jaeger-client-go/config" 实现快速构造 opentracing.Tracer
	//> 这里只是展开了其中的一些细节

	logSpans := true
	poolSpans := true
	bufferFlushInterval := 2 * time.Second
	logger := jaeger.StdLogger
	tracerMetrics := jaeger.NewMetrics(nil /*opts.metrics*/, nil)

	samplerMaker := func(types, parms string) jaeger.Sampler {
		var sampler jaeger.Sampler
		switch types {
		case jaeger.SamplerTypeProbabilistic:
			samplingRate, err := strconv.ParseFloat(parms, 32)
			if err == nil {
				sampler, _ = jaeger.NewProbabilisticSampler(samplingRate)
				break
			}
			fallthrough
		case jaeger.SamplerTypeConst:
			fallthrough
		default:
			constSample := true
			sampler = jaeger.NewConstSampler(constSample)
		}
		return sampler
	}
	sampler := samplerMaker(sampleType, sampleParm)

	reporterMaker := func() (jaeger.Reporter, error) {
		sender, err := jaeger.NewUDPTransport(agentAddr, 0)
		if err != nil {
			return nil, err
		}
		var reporter jaeger.Reporter
		reporterRemote := jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.QueueSize(0),
			jaeger.ReporterOptions.BufferFlushInterval(bufferFlushInterval),
			jaeger.ReporterOptions.Logger(logger),
			jaeger.ReporterOptions.Metrics(tracerMetrics))
		if logSpans && logger != nil {
			reporter = jaeger.NewCompositeReporter(jaeger.NewLoggingReporter(logger), reporterRemote)
		} else {
			reporter = reporterRemote
		}
		return reporter, nil
	}
	reporter, err := reporterMaker()
	if err != nil {
		return nil, nil
	}

	tracerOptions := []jaeger.TracerOption{
		jaeger.TracerOptions.PoolSpans(poolSpans),
		jaeger.TracerOptions.Logger(logger),
		jaeger.TracerOptions.Metrics(tracerMetrics),
	}
	tracer, closer := jaeger.NewTracer(svc, sampler, reporter, tracerOptions...)
	if tracer != nil {
		opentracing.SetGlobalTracer(tracer)
	}
	return tracer, closer
}

// JaegerUnaryServerInterceptor ...
func JaegerUnaryServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return jaegerUnaryServerInterceptor(tracer, jaegerGrpcFilterFunc)
}

// JaegerUnaryClientInterceptor ...
func JaegerUnaryClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return jaegerUnaryClientInterceptor(tracer, jaegerGrpcFilterFunc)
}

func jaegerUnaryServerInterceptor(tracer opentracing.Tracer, filter jaegerGrpcFilter) grpc.UnaryServerInterceptor {
	jaeger := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if filter(info.FullMethod) || nil == tracer {
			return handler(ctx, req)
		}

		currSvcAddr := ""
		startSpanOtions := func() []opentracing.StartSpanOption {
			var options []opentracing.StartSpanOption
			options = append(options, opentracing.Tag{Key: string(ext.Component), Value: "gRPC"})
			mdIn, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				mdIn = metadata.New(nil)
			} else {
				mdIn = mdIn.Copy()
			}
			ctx = metadata.NewOutgoingContext(ctx, mdIn)
			currSvcAddr = func() string {
				vals := mdIn.Get(tracingTargetAddr)
				if len(vals) > 0 {
					return vals[0]
				}
				return ""
			}()

			spanContext, err := tracer.Extract(opentracing.TextMap, textMapCarrier{mdIn})
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				xlog.Error("jaegerUnaryServerInterceptor extract from metadata err: %v", err)
			} else {
				options = append(options, opentracing.ChildOf(spanContext))
			}
			return options
		}

		//> TODO : trace caller
		//>      key = jaeger.TraceBaggageHeaderPrefix + "key" 才能在 SpanContext 中传播

		span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, info.FullMethod, startSpanOtions()...)
		defer span.Finish()
		if scTmp, ok := span.Context().(jaeger.SpanContext); ok {
			span.SetTag("SpanID", scTmp.SpanID().String())
			span.SetTag("ParentSpanID", scTmp.ParentID().String())
			span.SetTag(tracingTargetAddr, currSvcAddr)
		}

		injectSpanCtx := func() context.Context {
			mdOut, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				mdOut = metadata.New(nil)
			} else {
				mdOut = mdOut.Copy()
			}
			carrier := textMapCarrier{mdOut}
			err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
			if err != nil {
				span.LogFields(log.String("inject-error", err.Error()))
			}
			return metadata.NewOutgoingContext(ctx, mdOut)
		}
		ctx = injectSpanCtx()

		return handler(ctx, req)
	}
	return jaeger
}

func jaegerUnaryClientInterceptor(tracer opentracing.Tracer, filter jaegerGrpcFilter) grpc.UnaryClientInterceptor {
	jaeger := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if filter(method) || nil == tracer {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		mdOut, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			mdOut = metadata.New(nil)
		} else {
			mdOut = mdOut.Copy()
		}
		mdOut.Set(tracingTargetAddr, cc.Target())
		ctx = metadata.NewOutgoingContext(ctx, mdOut)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	return jaeger
}
