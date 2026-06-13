build:
	go build -o bin/mistah .

run:
	go run . $(ARGS)

test:
	go test ./...

clean-bin:
	rm -rf bin/
