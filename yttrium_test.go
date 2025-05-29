package yttrium

import (
	"context"
	"sync"
	"testing"
	"time"
)

// mockBuild implements AsyncBuild and simulates task execution.
type mockBuild struct {
	id        string
	callOrder *[]string
	delay     time.Duration
	lock      *sync.RWMutex
}

func (m *mockBuild) Setup(*Yttrium) AsyncBuild {
	return m
}

func (m *mockBuild) Run(_ *Yttrium, mu *sync.RWMutex, _ context.Context) {
	mu.Lock()
	defer mu.Unlock() // guarantee unlock even if panic happens

	time.Sleep(m.delay)

	*m.callOrder = append(*m.callOrder, m.id)
}

func TestRunTasksWithDAG(t *testing.T) {
	var order []string

	// Define builds
	a := &mockBuild{id: "A", callOrder: &order, delay: 10 * time.Millisecond}
	b := &mockBuild{id: "B", callOrder: &order, delay: 20 * time.Millisecond}
	c := &mockBuild{id: "C", callOrder: &order, delay: 5 * time.Millisecond}
	d := &mockBuild{id: "D", callOrder: &order, delay: 5 * time.Millisecond}

	// Create tasks: D depends on B and C; B depends on A
	taskA := NewTask(a)
	taskB := NewTask(b, taskA)
	taskC := NewTask(c)
	taskD := NewTask(d, taskB, taskC)

	yt := New()
	err := yt.RunTasks(context.Background(), taskD, taskB, taskC, taskA)
	if err != nil {
		t.Fatalf("RunTasks failed: %v", err)
	}

	// Verify topological order
	orderMap := make(map[string]int)
	for i, id := range order {
		orderMap[id] = i
	}

	if orderMap["A"] > orderMap["B"] {
		t.Error("Expected A before B")
	}
	if orderMap["B"] > orderMap["D"] {
		t.Error("Expected B before D")
	}
	if orderMap["C"] > orderMap["D"] {
		t.Error("Expected C before D")
	}
}
