package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/temporalio/samples-go/limit-workflow-tasks"
	"golang.org/x/sync/errgroup"

	"go.temporal.io/sdk/client"
)

type ClientSpec struct {
	TaskQueue         string
	NumberOfWorkflows int
	Concurrency       int
}

func parseFlags() (client.Options, ClientSpec) {
	set := flag.NewFlagSet("starter", flag.ContinueOnError)
	taskQueue := set.String("task-queue", "limit", "task queue name to isolate tests")
	numOfWorkflows := set.Int("num-workflows", 5, "number of unique workflows to be queried concurrently in the test")
	concurrency := set.Int("concurrency", 10, "concurrent clients executing workflows")

	clientOptions, err := limit.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		panic(err)
	}

	spec := ClientSpec{
		TaskQueue:         *taskQueue,
		NumberOfWorkflows: *numOfWorkflows,
		Concurrency:       *concurrency,
	}

	return clientOptions, spec
}

func main() {
	clientOptions, spec := parseFlags()
	fmt.Printf("Running starter with spec: %+v\n", spec)

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	bar := progressbar.Default(int64(spec.NumberOfWorkflows*spec.Concurrency), "🚀 executing workflows")
	then := time.Now()
	runs := &limit.Runs{}

	g := new(errgroup.Group)
	for i := 0; i < spec.Concurrency; i++ {
		g.Go(func() error {
			for i := 0; i < spec.NumberOfWorkflows; i++ {
				wo := client.StartWorkflowOptions{
					ID:        "workflow_limit_" + time.Now().String(),
					TaskQueue: spec.TaskQueue,
				}

				we, err := c.ExecuteWorkflow(context.Background(), wo, limit.PerformanceWorkflow)
				if err != nil {
					log.Fatalln("Unable to execute workflow", err)
				}
				bar.Add(1)
				runs.Add(we)
			}
			return nil
		})
	}
	g.Wait()

	time.Sleep(40 * time.Second)

	bar = progressbar.Default(int64(len(runs.Executions)), "🚀 getting workflow results")
	max := &limit.Max{}
	for _, run := range runs.Executions {
		var alloc int
		err = run.Get(context.Background(), &alloc)
		if err != nil {
			log.Fatalln("Unable to get workflow result", err)
		}
		max.Compare(alloc)
		bar.Add(1)
	}

	fmt.Printf("RPS: %.2f\n", float64(spec.NumberOfWorkflows*spec.Concurrency)/time.Since(then).Seconds())
	fmt.Printf("Max alloc = %v MiB\n", max.M)
}
