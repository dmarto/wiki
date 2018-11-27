.PHONY: build clean run build-run

wiki: main.go
	go get -t ./
	go build  ./

build: wiki

clean:
	rm -f wiki

run: build
	./wiki

build-run: build run
