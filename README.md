# apim-responser-tester

make all VERSION=v0.0.x
make push VERSION=v0.0.x

podman pull docker.io/golang:1.24.5
podman run -d --rm --name apim-response-tester -p 8080:8080 localhost/apim-response-tester:v0.0.3

