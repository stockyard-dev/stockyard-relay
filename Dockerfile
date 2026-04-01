FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go mod download && CGO_ENABLED=0 go build -o relay ./cmd/relay/

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/relay .
ENV PORT=9090 DATA_DIR=/data
EXPOSE 9090
CMD ["./relay"]
