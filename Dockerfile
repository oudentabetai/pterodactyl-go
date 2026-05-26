FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/pterodactyl-go .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /out/pterodactyl-go /app/pterodactyl-go

ENTRYPOINT ["/app/pterodactyl-go"]
