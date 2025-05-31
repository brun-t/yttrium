# Yttrium - The new next thing for build things

Yttrium is a "build engine" this means is 100% centered on build things

and that is just abstract enough for build anything but easily

# What does Build engine means?

An Build engine it just refer to a tool or in this case lib

capable of build anything but with enough commodities is not just an script

a example could be Docker this tool allow you to build things why is not an build engine?

well Docker is centered on Containers not building

# Okay so how I use this wonderful new thing?

Is really simple actually!

1.  Have go installed I don't think I need to explain this to you

    there are tons of tutorials for this online and that can

    explain it better than I could.

2.  Make init a Go project

    First make a folder into your current project where you want to use Yttrium

    and then init a go module here is the show of it!

    ```sh
    mkdir build
    cd build
    go mod build-MYPROJECTNAME # you also can do github.com/MyUser/MyProject/build
    ```

3.  Make your first Yttrium program!

okay so now make an main.go file yeah no trick no random file that need a cli just normal main.go file

so let's try an simple Build here,

NOTE:this doesn't build nothing is just an example of basic Yttrium apis aka Build and CommandRunner

main.go:

```go
package main

import (
    "fmt"
    "github.com/brun-t/yttrium"
)

type BasicBuild struct {
    file string
}

func (bb BasicBuild) Setup(yt *Yttrium) Build {
    return bb // not any init need
}

func (bb BasicBuild) Run(yt *Yttrium) error {
    cr := Yttrium.NewCommandRunner()

    result, err := cr.Exec("cat", bb.file)

    fmt.Printf("%s contents:%s", bb.file, string(result))

    return err
}

func main() {
    yt := Yttrium.New()

    bb := yt.Use(BasicBuild{file:"myfile.txt"})

    panic(yt.Run(bb))
}

```

and last step let's run it!

just do this in your shell

```sh
go mod tidy # this install Yttrium in your local project
go run . # this is the actual command that runs the app
```

