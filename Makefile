ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

build: build.vue build.go

build.vue:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd $(ROOT)/ui && npm install && npm run build
	@echo "\n"

build.go: 
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd $(ROOT)
	CGO_ENABLED=0 go build -o registry-webui main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main_linux_amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main_linux_arm64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o main_mac_amd64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o main_mac_arm64 main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o main.exe main.go
	@echo "\n"

test:
	go test -v ./... -cover

.PHONY: docker
docker: 
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	docker buildx build --file Dockerfile --platform linux/amd64,linux/arm64 --push -t seastart/registry-webui:latest .
	@echo "\n"
