package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/temporalio/samples-go/query-bench"
	"golang.org/x/sync/errgroup"

	"github.com/schollz/progressbar/v3"
	"go.temporal.io/sdk/client"
)

type ClientSpec struct {
	Iterations        int
	PayloadSize       int
	TaskQueue         string
	NumberOfWorkflows int
	QueryIntervalMs   int
}

func parseFlags() (client.Options, ClientSpec) {
	set := flag.NewFlagSet("query-worker", flag.ContinueOnError)
	iterations := set.Int("iterations", 200, "number of consecutive requests to be made from a client")
	payloadSize := set.Int("payload-size", 5, "size of the query result in bytes")
	taskQueue := set.String("task-queue", "query", "task queue name to isolate tests")
	numOfWorkflows := set.Int("num-workflows", 5, "number of unique workflows to be queried concurrently in the test")
	queryIntervalMs := set.Int("query-interval-ms", 250, "intervals between each query (milliseconds)")

	clientOptions, err := query.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		panic(err)
	}

	spec := ClientSpec{
		Iterations:        *iterations,
		PayloadSize:       *payloadSize,
		TaskQueue:         *taskQueue,
		NumberOfWorkflows: *numOfWorkflows,
		QueryIntervalMs:   *queryIntervalMs,
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

	executions := []client.WorkflowRun{}
	for i := 0; i < spec.NumberOfWorkflows; i++ {
		workflowOptions := client.StartWorkflowOptions{
			ID:        "query_workflow_" + time.Now().String(),
			TaskQueue: spec.TaskQueue,
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, query.PerformanceWorkflow, spec.PayloadSize)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		executions = append(executions, we)
	}

	log.Println("💤 Sleeping for workflow hydration")
	time.Sleep(45 * time.Second)
	log.Println("💤 Hydration completed")

	ctx := context.Background()
	q, err := query.NewQuantile()
	if err != nil {
		log.Fatalln(err)
	}
	then := time.Now()
	err = ProcessParallelQueries(ctx, c, executions, q, spec)
	if err != nil {
		fmt.Println("❌ Tests completed with error", err)
	}

	fmt.Println("✅ All tests completed")
	fmt.Printf("Test spec: %+v\n", spec)
	fmt.Printf("Test executed in between: %s - %s\n", then.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	fmt.Printf("Time elapsed: %v\n", time.Since(then))
	fmt.Printf("Total queries: %v\n", spec.NumberOfWorkflows*spec.Iterations)
	fmt.Printf("RPS: %.2f\n", float64(spec.NumberOfWorkflows*spec.Iterations)/time.Since(then).Seconds())
	fmt.Println("Query latency:")
	fmt.Printf("P(50): %d ms\n", time.Duration(q.Q(0.5)).Milliseconds())
	fmt.Printf("P(90): %d ms\n", time.Duration(q.Q(0.9)).Milliseconds())
	fmt.Printf("P(99): %d ms\n", time.Duration(q.Q(0.99)).Milliseconds())
}

func ProcessParallelQueries(ctx context.Context, c client.Client, we []client.WorkflowRun, q *query.Quantile, spec ClientSpec) error {
	bar := progressbar.Default(int64(spec.Iterations*len(we)), "🚀 Sending queries")
	g, ctx := errgroup.WithContext(ctx)
	for _, execution := range we {
		g.Go(func() error {
			return ProcessQuery(ctx, c, execution, q, bar, spec.Iterations, spec.QueryIntervalMs)
		})
	}
	return g.Wait()
}

func ProcessQuery(ctx context.Context, c client.Client, we client.WorkflowRun, q *query.Quantile, bar *progressbar.ProgressBar, iterations, interval int) error {
	for i := 0; i < iterations; i++ {
		var res string
		then := time.Now()

		value, err := c.QueryWorkflow(ctx, we.GetID(), we.GetRunID(), "state")
		if err != nil {
			return err
		}
		err = value.Get(&res)
		if err != nil {
			return err
		}

		q.Add(time.Since(then))
		bar.Add(1)

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	fmt.Print("\n")

	return nil
}
