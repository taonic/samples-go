package zapadapter

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a workflow function which does some logging.
// Important note: workflow logger is replay aware and it won't log during replay.
func Workflow(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := WorkflowLogger(ctx)
	logger.Info("Logging from workflow", "uuid", uuid.New().String())

	var result interface{}
	err := workflow.ExecuteActivity(ctx, LoggingActivity, "test").Get(ctx, &result)
	if err != nil {
		logger.Error("LoggingActivity failed.", "Error", err)
		return err
	}

	logger.Info("Workflow completed.")
	return nil
}

func WorkflowLogger(ctx workflow.Context) log.Logger {
	newLogger := workflow.GetLogger(ctx)
	if ctx != nil {
		newLogger = newLogger.(log.WithLogger).With("TraceId1", uuid.New().String())
		newLogger = newLogger.(log.WithLogger).With("TraceId2", uuid.New().String())
	}
	return newLogger
}

func LoggingActivity(ctx context.Context, name string) error {
	time.Sleep(8 * time.Second)
	logger := activity.GetLogger(ctx)
	logger.Info("Executing LoggingActivity.", "name", name)
	logger.Debug("Debugging LoggingActivity.", "value", "important debug data")
	return nil
}

func LoggingErrorAcctivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Warn("Ignore next error message. It is just for demo purpose.")
	logger.Error("Unable to execute LoggingErrorAcctivity.", "error", errors.New("random error"))
	return nil
}
