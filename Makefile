.PHONY: build build-run clean

build: clean
	go get -t ./
	go build  ./

build-run: build
	./wiki

clean:
	rm -f wiki
