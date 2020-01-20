# gomakecb

Go Make Crossplatform Builder.

`gomakecb` is the tool for a quick and easy cross-platform building of Golang projects. By the one command a source code could be compile for multiple architectures/operation systems. `gomakecb` implement the fully-functional wrapper for 'go' and 'make' commands, including a support pass an environmental variables and commandline's flags. The format of OS/ARCH is fully compatible with an output of `go tool dist list` command.
Also `gomakecb` could be used for cases when need an authomatize building's process, but some restrictions doesn't allow to get environments variables (e.g. building to a remote host).

## Quick start

Just download binary for your OS/ARCH and put that to an appropriate directory (for UNIX-like OSes downloaded binary must be an executable). After a preparation completed, follow to a directory, which containts a project. Then you can build a project from sources using `gomakecb`.
In a simplest case, when a build process imply only a call `go build...`,  using of `gomakecb` is an easy:
```
`which gomakecb` -t "go" -m="build" -osarch="linux/amd64,windows/amd64" -p="-o bin/\$GOOS/\$GOARCH/app -v app.go"
```
where `app` is a name of the output file and `app.go` is the name of source file. After the process of building will be finish, binaries will be stored in a directory tree (in `bin/`).  That's all :)
When the process of building is a more complicated and `make` utility is using, the step, after an installation, should including to edit `Makefile`. That's is an easy. Just add two lines (see below):
```
export GOOS
export GOARCH
```
and run:
```
./gomakecb -t "make" -f="Makefile" -m="build" -osarch="linux/amd64,windows/amd64"
```

