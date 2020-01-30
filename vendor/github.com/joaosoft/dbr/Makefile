env:
	docker-compose up -d home.postgres
run:
	go run ./main.go

build:
	go build .

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*