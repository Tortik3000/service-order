FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-order ./cmd/service-order

FROM gcr.io/distroless/base-debian12 AS service
COPY --from=builder /app/service-order /service-order
ENTRYPOINT ["/service-order"]
