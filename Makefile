export GOOS
export GOARCH
MAJOR_VERSION=1.0
BUILDNUM=21
BRANCH=prod
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gomakecb
PROJECT_FILES=gomakecb.go
LDFLAGS=-ldflags="-s -w -X main.version=$(MAJOR_VERSION) -X main.branch=$(BRANCH) -X main.buildnum=$(BUILDNUM) -X main.builddate=`date +%Y%m%d.%H%M%S` -X main.buildtime=`date +%s`"   
all: check test build
build: 
	$(GOBUILD) $(LDFLAGS) -o bin/$(GOOS)/$(GOARCH)/$(BINARY_NAME) -v $(PROJECT_FILES) 
test: 
	$(GOTEST) -v $(PROJECT_FILES)
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(PROJECT_FILES)
	./$(BINARY_NAME)
check:
	@[ "${GOOS}" ] || ( echo "GOOS is not set"; exit 1 )
	@[ "${GOARCH}" ] || ( echo "GOARCH is not set"; exit 1 )
