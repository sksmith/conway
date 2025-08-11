# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o conway-basic ./examples/basic
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o conway-advanced ./examples/advanced

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/conway-basic .
COPY --from=builder /app/conway-advanced .

# Make binaries executable
RUN chmod +x conway-basic conway-advanced

CMD ["./conway-basic"]