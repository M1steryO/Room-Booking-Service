FROM golang:1.23 AS builder
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY ../../Downloads/room-booking-service .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/room-booking-usecase ./cmd/api

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /out/room-booking-service /app/room-booking-service

EXPOSE 8080
ENTRYPOINT ["/app/room-booking-service"]
CMD ["serve"]
