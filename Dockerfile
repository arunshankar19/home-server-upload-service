FROM golang:1.25.2 AS builder
WORKDIR /upload_service
ENV GOPRIVATE=github.com/arunshankar19
COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=secret,id=netrc \
    cp /run/secrets/netrc ~/.netrc && \
    chmod 600 ~/.netrc && \
    go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o upload-service ./app/

FROM alpine:3.22
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /upload_service
COPY --chown=appuser:appgroup --from=builder /upload_service/upload-service .
USER appuser
CMD ["./upload-service"]
