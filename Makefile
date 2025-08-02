# Makefile for building apim-response-tester images for arm64 and amd64

IMAGE_NAME=apim-response-tester
VERSION=v0.0.3

.PHONY: all arm64 amd64

all: arm64 amd64
push: parm64 pamd64

arm64:
	podman build \
		--tag apim-response-tester-arm64:$(VERSION) \
		--build-arg ARCH=arm64 \
		.

amd64:
	podman build \
		--tag apim-response-tester-amd64:$(VERSION) \
		--build-arg ARCH=amd64 \
		.

parm64:
	podman push \
		localhost/apim-response-tester-arm64:$(VERSION) \
		docker.io/maclighiche/apim-response-tester-arm64:$(VERSION)

pamd64:
	podman push \
		localhost/apim-response-tester-amd64:$(VERSION) \
		docker.io/maclighiche/apim-response-tester-amd64:$(VERSION)

