package limit

import (
	"fmt"
	"runtime"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

func PerformanceWorkflow(ctx workflow.Context) (int, error) {
	logger := workflow.GetLogger(ctx)
	logger = log.With(logger, "uuid", string(make([]byte, 1024*16)))
	logger.Info("PerformanceWorkflow completed")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(bToMb(m.Alloc)), nil
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
