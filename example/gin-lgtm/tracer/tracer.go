package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Config struct {
	ServiceName string  `json:",optional"`
	Endpoint    string  `json:",optional"`
	Sampler     float64 `json:",default=1.0"`
	Batcher     string  `json:",default=otlphttp,options=zipkin|otlphttp|otlpgrpc"`
}

const (
	kindZipKin   = "zipkin"
	kindOtlphttp = "otlphttp"
	kindOtlpgrpc = "otlpgrpc"
)

var (
	tp      *tracesdk.TracerProvider
	_tracer trace.Tracer
)

type Span struct {
	span    trace.Span
	spanCtx context.Context
}

func InitTracer(c Config) error {
	var err error

	err = tracerProvider(c)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err.Error()))
		return err
	}

	otel.SetTracerProvider(tp)

	return nil
}

// tracerProvider is 返回一个openTelemetry TraceProvider，这里用的是jaeger
func tracerProvider(c Config) error {
	fmt.Println("init traceProvider")

	exp, err := createExporter(c)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err.Error()))
		return err
	}

	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.ServiceName),
		)),
	)

	_tracer = tp.Tracer(c.ServiceName + "-tracer")

	return nil
}

func createExporter(c Config) (tracesdk.SpanExporter, error) {
	// Just support jaeger and zipkin now, more for later
	switch c.Batcher {
	case kindZipKin:
		return zipkin.New(c.Endpoint)
	case kindOtlphttp:
		opts := []otlptracehttp.Option{
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(c.Endpoint),
		}
		return otlptracehttp.New(context.Background(), opts...)

	case kindOtlpgrpc:

		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(c.Endpoint),
		}
		return otlptracegrpc.New(context.Background(), opts...)

	default:
		return nil, fmt.Errorf("unknown exporter: %s", c.Batcher)
	}
}

// NewSpan is 初始化新span
func NewSpan(ctx context.Context, spanName string) Span {
	spanCtx, span := _tracer.Start(ctx, spanName)
	return Span{
		span:    span,
		spanCtx: spanCtx,
	}
}

func (S *Span) TraceID() string {
	return S.span.SpanContext().TraceID().String()
}

func (S *Span) SpanID() string {
	return S.span.SpanContext().SpanID().String()
}

func (S *Span) End() {
	S.span.End()
}

func (S *Span) Span() trace.Span {
	return S.span
}

func (S *Span) SpanContext() context.Context {
	return S.spanCtx
}

func (S *Span) NewChildSpan(spanName string) Span {
	return NewSpan(S.spanCtx, spanName)
}

const (
	HeaderTraceIDKey      = "x-trace-id"
	HeaderSpanIDKey       = "x-span-id"
	HeaderTraceContextKey = "x-trace-ctx"
)

func (S *Span) Inject(req *http.Request) {
	// header写入trace-id和span-id
	req.Header.Set(HeaderTraceIDKey, S.span.SpanContext().TraceID().String())
	req.Header.Set(HeaderSpanIDKey, S.span.SpanContext().SpanID().String())

	p := otel.GetTextMapPropagator()
	p.Inject(S.spanCtx, propagation.HeaderCarrier(req.Header))

	return
}

func (S *Span) InjectByBaggage(req *http.Request) {
	// 使用baggage写入trace id和span id
	p := propagation.Baggage{}

	traceMember, _ := baggage.NewMember(HeaderTraceIDKey, S.span.SpanContext().TraceID().String())
	spanMember, _ := baggage.NewMember(HeaderSpanIDKey, S.span.SpanContext().SpanID().String())

	b, _ := baggage.New(traceMember, spanMember)

	ctxBaggage := baggage.ContextWithBaggage(S.spanCtx, b)

	p.Inject(ctxBaggage, propagation.HeaderCarrier(req.Header))
}

func Extract(req *http.Request, spanName string) Span {
	var propagator = otel.GetTextMapPropagator()
	pctx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

	traceIDInHeader := req.Header.Get(HeaderTraceIDKey)
	spanIDInHeader := req.Header.Get(HeaderSpanIDKey)

	traceID, _ := trace.TraceIDFromHex(traceIDInHeader)
	spanID, _ := trace.SpanIDFromHex(spanIDInHeader)

	spanTempCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled, //这个没写，是不会记录的
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	spanRemoteCtx := trace.ContextWithRemoteSpanContext(pctx, spanTempCtx)

	spanCtx, span := _tracer.Start(spanRemoteCtx, spanName)

	return Span{span: span, spanCtx: spanCtx}
}

func ExtractByBaggage(req *http.Request, spanName string) Span {
	var propagator = propagation.TextMapPropagator(propagation.Baggage{})
	pctx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

	bag := baggage.FromContext(pctx)

	traceID, _ := trace.TraceIDFromHex(bag.Member(HeaderTraceIDKey).Value())
	spanID, _ := trace.SpanIDFromHex(bag.Member(HeaderSpanIDKey).Value())

	spanTempCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled, //这个没写，是不会记录的
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	spanRemoteCtx := trace.ContextWithRemoteSpanContext(pctx, spanTempCtx)

	spanCtx, span := _tracer.Start(spanRemoteCtx, spanName)

	return Span{span: span, spanCtx: spanCtx}
}
