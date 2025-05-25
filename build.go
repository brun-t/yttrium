package yttrium

import (
	"context"
	"sync"
)

/*
This interface define what a struct need for can be used as a runnable in Yttrium
*/
type Build interface {
	Setup(*Yttrium) Build
	Run(*Yttrium) error
}

/*
This interface define what a struct need for can be used as a runnable in Yttrium
*/
type AsyncBuild interface {
	Setup(*Yttrium) AsyncBuild
	//sync.RWMutex is for manage read/writing, context.Context is for handle cancellation signals
	Run(*Yttrium, *sync.RWMutex, context.Context)
}
