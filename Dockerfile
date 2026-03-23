FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /out/room-booking-service ./cmd/room-booking-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /out/room-booking-service /app/room-booking-service

EXPOSE 8080

ENTRYPOINT ["/app/room-booking-service"]
CMD ["serve"]
