package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/temporalio/samples-go/zapadapter"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"

	prom "github.com/prometheus/client_golang/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	setupProfiler()
	defer profiler.Stop()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6065", nil))
	}()

	// Server for pprof
	c, err := client.Dial(client.Options{
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:8077",
			TimerType:     "histogram",
		})),
		// ZapAdapter implements log.Logger interface and can be passed
		// to the client constructor using client using client.Options.
		Logger: zapadapter.NewZapAdapter(
			NewZapLogger()),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "zap-logger-5", worker.Options{
		MaxConcurrentWorkflowTaskPollers: 10,
	})

	w.RegisterWorkflow(zapadapter.Workflow)
	w.RegisterActivity(zapadapter.LoggingActivity)
	w.RegisterActivity(zapadapter.LoggingErrorAcctivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func NewZapLogger() *zap.Logger {
	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Sampling:         nil, // consider exposing this to config for our external customer
		Encoding:         "json",
		EncoderConfig:    encodeConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()

	// Or simple
	// logger, err := zap.NewDevelopment()
	// can be used instead of the code above.

	if err != nil {
		log.Fatalln("Unable to create zap logger")
	}
	return logger
}

func setupProfiler() {
	err := profiler.Start(
		profiler.WithService("zap-tests"),
		profiler.WithEnv("dev"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
