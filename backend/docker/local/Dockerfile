FROM golang:1.26-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@v1.63.4

COPY go.mod go.sum ./
RUN go mod download

# Run air for hot reloading
CMD ["air"]
