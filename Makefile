.PHONY: build clean run

wiki: main.go
	go get -t ./
	go build  ./

build: wiki

clean:
	rm -f wiki

run: build
	./wiki
