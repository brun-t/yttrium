package yttrium

/*
This interface define what a struct need for can be used as a runnable in Yttrium
*/
type Build interface {
	Setup(*Yttrium) *Build
	Run(*Yttrium) error
}
