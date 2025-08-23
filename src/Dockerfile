# ---------- Builder Stage ----------
FROM golang:1.25 as builder


WORKDIR /app


COPY src/go.mod src/go.sum ./


RUN go mod download


COPY src/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o abrarvan_challenge ./cmd

# ---------- Runtime Stage ----------
FROM alpine:3.22 as runtime


RUN apk --no-cache add \
    ca-certificates \
    wget \
    curl \
    tzdata \
    && update-ca-certificates


RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup


WORKDIR /app


COPY --from=builder --chown=appuser:appgroup /app/abrarvan_challenge .


COPY --from=builder --chown=appuser:appgroup /app/config ./config

RUN chmod +x ./abrarvan_challenge

RUN mkdir -p logs && chown -R appuser:appgroup logs


USER appuser

EXPOSE 5005


CMD ["./abrarvan_challenge"]
