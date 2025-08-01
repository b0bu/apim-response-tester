# apim-responser-tester

podman pull docker.io/golang:1.24.5
podman build --tag apim-response-tester:v0.0.2 .
podman run -d --rm --name apim-response-tester -p 8080:8080 localhost/apim-response-tester:v0.0.3

