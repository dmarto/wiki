.PHONY: dev build clean

dev: build
	./wiki

build: clean
	go get -t ./
	go build ./

test:
	go test ./

clean:
	rm -rf wiki
