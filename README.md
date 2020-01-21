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

## Environment variables, commandline arguments

In some cases, may be necessary to passing a special environment variables when calling 'go' or 'make' commands. This may be required when a compilation process will be done to a remote host. To pass of an additional variables for 'make' or 'go' should be using commandline switch '-e'. E.g.:
```
./gomakecb -t "make" -f="Makefile" -m="build" -osarch="linux/amd64,windows/amd64" -e="BUILDMODE=prod"
```
Another scenario implies to overwrite ща environment variables inherited for a user. For this case has been implemented the switch '-eow' (overwrite of environments variables). If `-eow` engaged, all environment variables will be replacted to values which were passed via '-e' switch. E.g.:
```
./gomakecb -t "go" -osarch="all" -m="build" -p="-ldflags='-s -w -X main.version=0.1 -X main.builddate=`date -u +%Y%m%d.%H%M%S`' -o bin/\$GOOS/\$GOARCH/app -v app.go" -e="HOME=/tmp,GOCACHE=/tmp,PATH=/sbin:/usr/sbin:/usr/local/sbin:/usr/local/bin:/usr/bin:/bin:/usr/local/go/bin" -eow -d
```
Сommandline arguments aren't require a detailed description, unlike an environment variables. For pass them just set  the `-p` switch and assign required values. E.g.:
```
./gomakecb -t "make" -f="Makefile" -m="build" -osarch="linux/amd64,windows/amd64" -p="TEST=true"
```

## Build-in help

```
Usage of bin/linux/amd64/gomakecb:
  -d    Debug output.
  -e string
        Environment variables.
  -eow
        Overwriting of environment variables.
  -f string
        Path to Makefile (only if -t 'make').
  -list
        Print the list of supported GOOS/GOARCH.
  -m string
        Build mode (e.g. 'build' or 'clear'). (default "build")
  -osarch string
        Set GOOS/GOARCH. Use 'all' for build for all OS/ARCHs. (default "linux/amd64")
  -p string
        Another parameters for 'make'/'go' which should be passed.
  -s    Perform a simulate mode.
  -t string
        Build tool: 'make' | 'go'.
  -timeout string
        Maximum timeout execution of 'make'/'go'. (default "24h")
  -v    Show version.

Examples:
bin/linux/amd64/gomakecb -t "make" -f="Makefile" -m="build" -osarch="linux/amd64,windows/amd64"
bin/linux/amd64/gomakecb -t "go" -m="build" -osarch="linux/amd64,windows/amd64" -p="-o bin/\$GOOS/\$GOARCH/app -v app.go"
```

## Supported OS/ARCHs

```
- List of supported GOOS/GOARCH:
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/386
darwin/amd64
darwin/arm
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/s390x
nacl/386
nacl/amd64p32
nacl/arm
netbsd/386
netbsd/amd64
netbsd/arm
openbsd/386
openbsd/amd64
openbsd/arm
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
windows/386
windows/amd64
windows/arm
- Total: 41
```
