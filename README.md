# Yttrium a Build engine for the new era

Yttrium is a Build engine that is used directly as a Go lib

that is mostly adapted for C/C++ but it can build anything else

So first you would be asking what is a "Build engine"?

## Build engine

This refer to a lib, framework or app that can be used for build anything in a simple way for example Docker!

Docker is not technically a pure a Build engine but is very close!

So in Docker you can describe instructions that build things

The thing is that docker is not focused on Build is focused on Containers and a "Build engine" is centered in building things

## How use Yttrium

You have your project so you first need to have Go installed

Okay so now you setup a basic Go project, here is a guide step to step of how do so if you aren't familiar to Go

The standard for make a Yttrium project is

```sh
mkdir build
cd build
go mod init build-YOUR-PROJECT-NAME-GOES-HERE
```

And now you make a new file called

main.go:

with this code for example

```go
package main

import (
    "github.com/brun-t/Yttrium"
)

func main() {
    yt := Yttrium.New()

    exe := yt.Executable("my-app")

    exe.AddSources("src/main.cpp")

    exe.AddIncludes("includes", "thirdparty/STB")

    yt.Run(exe)
}
```

then you do this

```sh
go mod tidy

# and now just exec it!

go run .
```

And done! you learned basic Yttrium setup
