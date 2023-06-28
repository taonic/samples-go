package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/mtls"
	query "github.com/temporalio/samples-go/query-bench"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, _, err := mtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	worker.EnableVerboseLogging(true)

	w := worker.New(c, "query", worker.Options{MaxConcurrentWorkflowTaskPollers: 50})

	w.RegisterWorkflow(query.QueryWorkflow)
	w.RegisterActivity(query.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
