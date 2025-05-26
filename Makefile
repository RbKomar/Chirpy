build:
	go build -o bin/main .

run: build
	./bin/main

clean:
	rm -f bin/main