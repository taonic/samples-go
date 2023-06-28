package query

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

const (
	numActivities    = 500 // this should produce ~1500 events
	sleepMinutes     = 60
	WorkflowIDPrefix = "query_workflow_"
)

// Workflow is to demo how to setup query handler
func QueryWorkflow(ctx workflow.Context) error {
	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (string, error) {
		return "hello", nil
	})

	logger := workflow.GetLogger(ctx)
	logger.Info("QueryWorkflow started")

	var result string
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for i := 1; i < numActivities; i++ {
		logger.Info("Starting activity 1")
		err = workflow.ExecuteActivity(ctx, Activity, "Hello").Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return err
		}
	}

	logger.Info(fmt.Sprintf("Putting workflow to sleep for %s minutes", sleepMinutes))
	_ = workflow.Sleep(ctx, sleepMinutes*time.Minute)

	logger.Info("Workflow completed")
	return nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}
