package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	limit "github.com/temporalio/samples-go/limit-workflow-tasks"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"

	"net/http"
	_ "net/http/pprof"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type WorkerSpec struct {
	MaxConcurrentWorkflowTaskPollers       int
	MaxConcurrentWorkflowTaskExecutionSize int
	TaskQueue                              string
}

func parseFlags() (client.Options, WorkerSpec) {
	// Server for pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	set := flag.NewFlagSet("worker", flag.ContinueOnError)
	maxConcurrentWorkflowTaskPollers := set.Int("max-pollers", 2, "Concurrent workflow task pollers")
	maxConcurrentWorkflowTaskExecutionSize := set.Int("max-exec-slots", 10000, "Concurrent workflow execution slots")
	taskQueue := set.String("task-queue", "limit", "task queue name to isolate tests")

	clientOptions, err := limit.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		panic(err)
	}

	workerSpec := WorkerSpec{
		MaxConcurrentWorkflowTaskPollers:       *maxConcurrentWorkflowTaskPollers,
		MaxConcurrentWorkflowTaskExecutionSize: *maxConcurrentWorkflowTaskExecutionSize,
		TaskQueue:                              *taskQueue,
	}

	return clientOptions, workerSpec
}

func main() {
	setupProfiler()
	defer profiler.Stop()

	clientOptions, spec := parseFlags()
	fmt.Printf("Running worker with spec: %+v\n", spec)

	clientOptions.Logger = limit.NewZapAdapter(limit.NewZapLogger())

	clientOptions.MetricsHandler = sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
		ListenAddress: "0.0.0.0:9090",
		TimerType:     "histogram",
	}))

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	worker.EnableVerboseLogging(true)

	w := worker.New(c, spec.TaskQueue, worker.Options{
		MaxConcurrentWorkflowTaskPollers:       spec.MaxConcurrentWorkflowTaskPollers,
		MaxConcurrentWorkflowTaskExecutionSize: spec.MaxConcurrentWorkflowTaskExecutionSize,
		WorkerActivitiesPerSecond:              1,
		TaskQueueActivitiesPerSecond:           1,
	})

	w.RegisterWorkflow(limit.PerformanceWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func setupProfiler() {
	err := profiler.Start(
		profiler.WithService("limit-tests"),
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
