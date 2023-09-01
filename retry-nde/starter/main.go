package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	"github.com/google/uuid"
	retrynde "github.com/temporalio/samples-go/retry-nde"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:          "hello_world_" + uuid.New().String(),
		TaskQueue:   "hello-world",
		RetryPolicy: &temporal.RetryPolicy{}, // use default retry policy
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, retrynde.Workflow, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
