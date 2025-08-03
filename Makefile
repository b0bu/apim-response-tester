# Makefile for building $(IMAGE_NAME) images for arm64 and amd64

IMAGE_NAME=apim-response-tester

.PHONY: all server-arm64 server-amd64 client

all: server-arm64 server-amd64 client
push: parm64 pamd64


#client:
#	GOOS=linux GOARCH=arm64 go build -o bin/api-response-tester-arm64 ./client
#	GOOS=linux GOARCH=amd64 go build -o bin/api-response-tester-amd64 ./client


server-arm64:
	podman build \
		--platform linux/arm64 \
		--tag $(IMAGE_NAME)-arm64:$(VERSION) \
		--build-arg ARCH=arm64 \
		.

server-amd64:
	podman build \
		--platform linux/amd64 \
		--tag $(IMAGE_NAME)-amd64:$(VERSION) \
		--build-arg ARCH=amd64 \
		.

parm64:
	podman push \
		localhost/$(IMAGE_NAME)-arm64:$(VERSION) \
		docker.io/maclighiche/$(IMAGE_NAME)-arm64:$(VERSION)

pamd64:
	podman push \
		localhost/$(IMAGE_NAME)-amd64:$(VERSION) \
		docker.io/maclighiche/$(IMAGE_NAME)-amd64:$(VERSION)

