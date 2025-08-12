FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /block-synchronizer cmd/block-synchronizer/main.go
RUN go build -o /event-processor cmd/event-processor/main.go
RUN go build -o /balance-api cmd/balance-api/main.go
