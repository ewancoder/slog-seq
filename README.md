# slog-seq

**slog-seq** is a library for sending logs to a [Seq](https://datalust.co/seq) server, as a handler for Go's structured logging [slog](https://go.dev/blog/slog).

It also supports some trace functionality.

## Installation

```bash
go get github.com/sokkalf/slog-seq
```

## Quick start

It's pretty easy to get going.

```go
seqLogger, handler := slogseq.NewSeqLogger(
	"http://your-seq-server/ingest/clef",
	"your-api-key",
	50,            // batchSize
	2*time.Second, // flushInterval
	nil,           // opts go here, if there are any
)
defer handler.Close()

slog.SetDefault(seqLogger)
slog.Info("Hello, world!")
```

You can set some options, here are some examples:

```go
opts := &slog.HandlerOptions{
	Level:     slog.LevelInfo,  // minimum log level
	AddSource: true,            // show source file, line and function in log
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "password" {
			// mask passwords
			a.Value = slog.StringValue("*****")
		}
		if a.Key == slog.SourceKey {
			// Replace the full path with just the file name
			s := a.Value.Any().(*slog.Source)
			s.File = path.Base(s.File)
		}
		return a
	},
}
```

If you need to disable TLS certificate verification, you can do so by calling `handler.DisableTLSVerification()`.

## Traces

```go
spanProcessor := trace.NewSimpleSpanProcessor(&slogseq.LoggingSpanProcessor{Handler: handler})
tp := trace.NewTracerProvider(trace.WithSpanProcessor(spanProcessor), trace.WithSampler(trace.AlwaysSample()))
tracer := tp.Tracer("example-tracer")
ctx := context.Background()
spanCtx, span := tracer.Start(ctx, "operation")
span.AddEvent("Starting work")
time.Sleep(500 * time.Millisecond)
slog.InfoContext(spanCtx, "This is a span log message", "key", "value")
spanCtx, subSpan := tracer.Start(spanCtx, "sub operation")
subSpan.AddEvent("Sub operation started")
time.Sleep(500 * time.Millisecond)
subSpan.AddEvent("Sub operation completed", tr.WithAttributes(attribute.String("key", "value")))
subSpan.End()
span.AddEvent("Work done")
slog.InfoContext(spanCtx, "All done!")
span.End()
```

![Seq with traces](../master/doc/seq_screenshot.png)

## License

MIT
