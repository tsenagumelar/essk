FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY services/backend/go.mod services/backend/go.sum ./
RUN go mod download

COPY services/backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/essk-server ./services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/essk ./cmd/essk

FROM alpine:3.21

RUN addgroup -S essk && adduser -S essk -G essk

WORKDIR /app

COPY --from=builder /out/essk-server /usr/local/bin/essk-server
COPY --from=builder /out/essk /usr/local/bin/essk
COPY services/backend/migrations ./migrations
COPY services/backend/shared ./shared

USER essk

EXPOSE 18080 19100

CMD ["essk-server"]
