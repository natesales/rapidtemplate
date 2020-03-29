all: install build

install:
	go get -u github.com/gomarkdown/markdown
	go get -u github.com/fsnotify/fsnotify

build:
	go clean
	go build