build:
	go install -v .

test:
	go test -v ./...

image:
	docker build -t cirocosta/slirunner .
