FROM golang:1.23-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 appuser && \
    adduser -D -s /bin/sh -u 1001 -G appuser appuser

WORKDIR /root/

COPY --from=builder /build/main .

RUN chown -R appuser:appuser /root/

# Use an unprivileged user
USER appuser

EXPOSE 8080

# Command to run the executable
CMD ["./main"] 