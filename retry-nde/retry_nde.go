package helloworld

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var result string

	random := rand.Float64()
	// There is a 50% chance to run into non-deterministic error
	if random > 0.5 {
		err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return "", err
		}
	}

	fmt.Println("====> Restart worker now.")
	workflow.Sleep(ctx, 5*time.Second)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}
