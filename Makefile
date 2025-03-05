build:
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./... -count=1

sqlcgen:
	docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate
