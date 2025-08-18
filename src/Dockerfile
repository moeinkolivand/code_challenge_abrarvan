# ---------- Builder Stage ----------
FROM golang:1.25 as builder
WORKDIR /app

# Copy go.mod / go.sum first and download deps
COPY src/go.mod src/go.sum ./

RUN go mod download

# Copy the full source
COPY src/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o abrarvan_challenge ./cmd

# ---------- Runtime Stage ----------
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/abrarvan_challenge .

COPY --from=builder /app/config ./config

ENTRYPOINT ["./abrarvan_challenge"]
