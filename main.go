package main

import (
	"context"
	"io"
	"os"

	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type logOrigin struct {
	File struct {
		Name     string `json:"name"`
		Line     int    `json:"line"`
		Function string `json:"function"`
	} `json:"file"`
}

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(h.addTraceFromContext(ctx)...)
	return h.Handler.Handle(ctx, r)
}

func (h ContextHandler) addTraceFromContext(ctx context.Context) (as []slog.Attr) {
	if ctx == nil {
		return
	}
	span := trace.SpanContextFromContext(ctx)
	traceID := span.TraceID().String()
	spanID := span.SpanID().String()
	traceGroup := slog.Group("trace", slog.String("id", traceID))
	spanGroup := slog.Group("span", slog.String("id", spanID))
	as = append(as, traceGroup)
	as = append(as, spanGroup)
	return
}

func getJsonHandler(w io.Writer) *slog.JSONHandler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "@timestamp"
			}
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}
			if a.Key == slog.LevelKey {
				a.Key = "log.level"
			}
			if a.Key == slog.SourceKey {
				a.Key = "log.origin"
				source := a.Value.Any().(*slog.Source)
				var file logOrigin
				file.File.Function = source.Function
				file.File.Name = source.File
				file.File.Line = source.Line
				a.Value = slog.AnyValue(file)
			}
			return a
		},
	})
}

func GetLogger() *slog.Logger {
	jsonHandler := getJsonHandler(os.Stdout)
	ctxHandler := ContextHandler{jsonHandler}
	return slog.New(ctxHandler)
}
