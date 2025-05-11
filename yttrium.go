package yttrium

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
func (yt *Yttrium) Use(b *Build) *Build {
	return (*b).Setup(yt)
}

// It runs a build
func (yt *Yttrium) Run(b *Build) error {
	return (*b).Run(yt)
}
