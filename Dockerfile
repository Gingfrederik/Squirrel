FROM golang:1.11.5 AS builder
RUN mkdir /fileserver
WORKDIR /fileserver
ENV GO111MODULE=on 

COPY go.mod . 
COPY go.sum .

RUN go mod download
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api

FROM alpine:latest
RUN mkdir -p /fileserver
COPY --from=builder /fileserver/api /fileserver/api
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /fileserver
USER appuser

WORKDIR /fileserver
CMD ./api
