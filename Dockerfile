FROM golang:1.20-alpine as builder
RUN apk --no-cache add git make build-base

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o /bin/relay ./cmd/relay

# --- Execution Stage

FROM alpine:latest

COPY --from=builder /bin/relay /bin/

EXPOSE 80
ENTRYPOINT ["/bin/relay"]
