package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	query "github.com/temporalio/samples-go/query-bench"

	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type WorkerSpec struct {
	MaxConcurrentWorkflowTaskPollers int
	TaskQueue                        string
}

func parseFlags() (client.Options, WorkerSpec) {
	set := flag.NewFlagSet("query-worker", flag.ContinueOnError)
	maxConcurrentWorkflowTaskPollers := set.Int("pollers", 10, "Concurrent workflow task pollers")
	taskQueue := set.String("task-queue", "query", "task queue name to isolate tests")

	clientOptions, err := query.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		panic(err)
	}

	workerSpec := WorkerSpec{
		MaxConcurrentWorkflowTaskPollers: *maxConcurrentWorkflowTaskPollers,
		TaskQueue:                        *taskQueue,
	}

	return clientOptions, workerSpec
}

func main() {
	setupProfiler()
	defer profiler.Stop()

	clientOptions, spec := parseFlags()
	fmt.Printf("Running worker with spec: %+v\n", spec)

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	worker.EnableVerboseLogging(true)

	w := worker.New(c, spec.TaskQueue, worker.Options{
		MaxConcurrentWorkflowTaskPollers: spec.MaxConcurrentWorkflowTaskPollers,
	})

	w.RegisterWorkflow(query.PerformanceWorkflow)
	w.RegisterActivity(query.MockActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func setupProfiler() {
	err := profiler.Start(
		profiler.WithService("query-tests"),
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
