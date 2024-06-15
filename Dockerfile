FROM golang:latest AS build

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/telegram-bot/main.go

FROM alpine:3
COPY --from=build build/main /bin/main
COPY dev.yaml .
COPY test.yaml .

ENTRYPOINT ["/bin/main"]
