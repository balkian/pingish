build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o pingish

docker-build:
	docker build -t balkian/pingish .

run: build
	./pingish -c 10 -host www.google.es
