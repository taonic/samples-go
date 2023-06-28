package main

import (
	"context"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/mtls"
	query "github.com/temporalio/samples-go/query-bench"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	clientOptions, wid, err := mtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        query.WorkflowIDPrefix + wid,
		TaskQueue: "query",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, query.QueryWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
