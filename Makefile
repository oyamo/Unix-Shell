all: build

build:
	go build -o osh

.PHONY: clean

clean:
	rm -rf ./build