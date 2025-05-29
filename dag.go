package yttrium

import (
	"fmt"

	"github.com/google/uuid"
)

type Task struct {
	ID    string
	Build AsyncBuild
	Deps  []*Task
}

func GenerateTaskID() string {
	id := uuid.New()
	return id.String()
}

func NewTask(build AsyncBuild, deps ...*Task) *Task {
	return &Task{
		ID:    GenerateTaskID(),
		Build: build,
		Deps:  deps,
	}
}

func TopoSort(tasks []*Task) ([]*Task, error) {
	inDegree := make(map[string]int)
	idToTask := make(map[string]*Task)

	// Collect all tasks and dependencies into the maps
	var collect func(t *Task)
	visited := make(map[string]bool)
	collect = func(t *Task) {
		if visited[t.ID] {
			return
		}
		visited[t.ID] = true
		idToTask[t.ID] = t
		if _, exists := inDegree[t.ID]; !exists {
			inDegree[t.ID] = 0
		}
		for _, dep := range t.Deps {
			collect(dep)
			inDegree[t.ID]++
		}
	}

	for _, task := range tasks {
		collect(task)
	}

	// Start with all nodes with in-degree 0
	queue := make([]*Task, 0)
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, idToTask[id])
		}
	}

	var sorted []*Task
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, t := range idToTask {
			for _, dep := range t.Deps {
				if dep.ID == current.ID {
					inDegree[t.ID]--
					if inDegree[t.ID] == 0 {
						queue = append(queue, t)
					}
				}
			}
		}
	}

	if len(sorted) != len(idToTask) {
		return nil, fmt.Errorf("cycle detected in DAG")
	}

	return sorted, nil
}
