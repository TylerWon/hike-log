# Build stage
# ---------------------
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0: produce a static binary that runs on Alpine without C deps
# -trimpath: strip local filesystem paths from the binary for more reproducible builds, less info leakage
# -ldflags="-s -w": strip the symbol table (-s) and DWARF debug info (-w) to shrink the binary
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

# Run stage
# ---------------------
FROM alpine:3.21

# Install the CA trust store so HTTPS calls (ex. to AWS S3) verify TLS
RUN apk add --no-cache ca-certificates

# Create user "appuser" with UID 1000 and no password
RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /out/server /app/server

USER appuser

EXPOSE 8080

CMD ["./server"]
