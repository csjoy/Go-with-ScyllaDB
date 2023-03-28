# Builder image
FROM docker.io/library/golang AS builder
COPY go.mod main.go /go/src/
WORKDIR /go/src
RUN go get -d -v ./...
RUN go build -v -o app ./...

# Runtime image
FROM docker.io/library/debian:stable-slim AS server
COPY --from=builder /go/src/app /
CMD ["/app"]