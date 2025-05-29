package yttrium

import (
	"context"
	"fmt"
	"sync"
)

/*
TODO:Find what we need here, I think with a blank struct this should work fine
TODO:because Yttrium is a runner so it don't to store data
*/

/*
Yttrium is the central Runner for can run all Build's
*/
type Yttrium struct{}

// Is the constructor for Yttrium struct aka it makes a new instance of Yttrium
func New() *Yttrium {
	return &Yttrium{}
}

// It setups a Build
func (yt *Yttrium) Use(b Build) Build {
	return b.Setup(yt)
}

// It runs a build
func (yt *Yttrium) Run(b Build) error {
	return b.Run(yt)
}

// It setups and AsyncBuild
// IMPORTANT: the Setup itself is not async this is just for a AsyncBuild
func (yt *Yttrium) AsyncUse(b AsyncBuild) AsyncBuild {
	return b.Setup(yt)
}

// it runs an AsyncBuild and accepts and context for cancellation signals, timeouts, etc
func (yt *Yttrium) AsyncRun(b AsyncBuild, ctx context.Context) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("panic: %v", r)
				}
			}
			wg.Done()
		}()
		b.Run(yt, &sync.RWMutex{}, ctx)
	}()

	wg.Wait()
	return
}

// It runs multiples instances of AsyncBuild
func (yt *Yttrium) AsyncRuns(ctx context.Context, b ...AsyncBuild) []error {
	var wg sync.WaitGroup
	var rwmutex sync.RWMutex
	errChan := make(chan error, len(b))

	wg.Add(len(b))
	for _, build := range b {
		go func(build AsyncBuild) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						errChan <- err
					} else {
						errChan <- fmt.Errorf("panic: %v", r)
					}
				}
				wg.Done()
			}()
			build.Run(yt, &rwmutex, ctx)
		}(build)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// FP rocks!
	return func(errChan chan error) []error {
		var errs []error
		for err := range errChan {
			errs = append(errs, err)
		}
		return errs
	}(errChan)
}

func (yt *Yttrium) RunTasks(ctx context.Context, t ...*Task) []error {
	sorted, err := TopoSort(t)

	if err != nil {
		return []error{err}
	}

	var mutex sync.RWMutex
	taskWaitGroups := make(map[string]*sync.WaitGroup)
	errChan := make(chan error, len(t))
	var wg sync.WaitGroup

	for _, task := range sorted {
		taskWaitGroups[task.ID] = &sync.WaitGroup{}
		taskWaitGroups[task.ID].Add(1)
	}

	wg.Add(len(sorted))
	for _, task := range sorted {
		go func(task *Task) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						errChan <- err
					} else {
						errChan <- fmt.Errorf("panic: %v", r)
					}
				}
				taskWaitGroups[task.ID].Done()
				wg.Done()
			}()

			// Wait for all dependencies to finish
			for _, dep := range task.Deps {
				taskWaitGroups[dep.ID].Wait()
			}

			// Run the build task
			task.Build.Run(yt, &mutex, ctx)
		}(task)
	}

	// Wait for all goroutines to finish and then close error channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	return errs
}
