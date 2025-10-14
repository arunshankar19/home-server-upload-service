FROM golang:1.25.2-alpine AS builder
WORKDIR /upload_service
RUN apk add --no-cache git
ENV GOPRIVATE=github.com/arunshankar19
COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o upload-service ./app/

FROM alpine:3.22
WORKDIR /upload_service
COPY --from=builder /upload_service/upload-service .
CMD ["./upload-service"]