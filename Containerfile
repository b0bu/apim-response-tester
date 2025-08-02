# syntax=docker/dockerfile:1

FROM golang:1.24.5

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -o /apim-response-tester

EXPOSE 8080

CMD ["/apim-response-tester"]
