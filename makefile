# Detect system OS.
ifeq ($(OS),Windows_NT)
    detected_OS := Windows
else
    detected_OS := $(shell sh -c 'uname -s 2>/dev/null || echo not')
endif

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -race
GOGET=$(GOCMD) get


ifeq ($(detected_OS),Windows)
	BINARY_NAME=getNovel.exe
	BINARY_RACE_NAME=getNovel_race.exe
else
	BINARY_NAME=getNovel
	BINARY_RACE_NAME=getNovel_race
endif


DIR_TEMP=temp/
DIR_RESULT=finish/

flags="-X 'main.goversion=`go version`' -X 'main.buildstamp=`date --rfc-3339=seconds`' -X main.githash=`git describe --always --long --abbrev=14`"

all: test build
build:	
	$(GOBUILD) -ldflags ${flags} -x   -v -o $(BINARY_NAME)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f engine.log
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f $(BINARY_RACE_NAME)
	rm -f $(BINARY_RACE_NAME).exe
	rm -rf $(DIR_TEMP)
	rm -rf $(DIR_RESULT)
run: build
	./$(BINARY_NAME) -version
race:
	$(GOBUILD)  -race -ldflags ${flags} -x   -v -o $(BINARY_RACE_NAME)
