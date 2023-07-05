package limit

import (
	"sync"

	"go.temporal.io/sdk/client"
)

type Runs struct {
	mu         sync.Mutex
	Executions []client.WorkflowRun
}

func (m *Runs) Add(exec client.WorkflowRun) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Executions = append(m.Executions, exec)
}
