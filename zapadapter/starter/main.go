package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/temporalio/samples-go/zapadapter"
	"go.temporal.io/sdk/client"
	"golang.org/x/sync/errgroup"
)

const concurrency = 10

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	g := new(errgroup.Group)
	for k := 0; k < concurrency; k++ {
		g.Go(func() error {
			for i := 0; i < 2000; i++ {
				workflowOptions := client.StartWorkflowOptions{
					ID:        "zap_logger_workflow_id" + uuid.New().String(),
					TaskQueue: "zap-logger-5",
				}

				we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, zapadapter.Workflow, "<param to log>")
				if err != nil {
					log.Fatalln("Unable to execute workflow", err)
				}

				log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
			}
			return nil
		})
	}

	g.Wait()
}
