# gomakecb

Go Make Crossplatform Builder.

## Description

`gomakecb` is the tool for an easy cross-platform building of Golang projects.

### Examples

#### -t 'go'

```
# cd examples/
../gomakecb -t "go" -osarch="all" -m="build" -p="-ldflags='-s -w -X main.version=0.1 -X main.builddate=`date -u +%Y%m%d.%H%M%S`' -o bin/\$GOOS/\$GOARCH/app -v app.go"
```

#### -t 'make'

```
# cd examples/
../gomakecb -t "make" -osarch="all" -m="build" 
```

