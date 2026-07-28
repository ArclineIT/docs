FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o docs-server ./cmd/serve

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/docs-server .
COPY content  ./content
COPY static   ./static
COPY templates ./templates

EXPOSE 8080

ENV PORT=8080 \
    CONTENT_DIR=/app/content \
    STATIC_DIR=/app/static \
    TEMPLATES_DIR=/app/templates \
    BILLING_DB=/data/billing/arcline-billing.db \
    DOCS_DB=/data/docs/arcline-docs.db \
    BILLING_URL=https://portal.arcline.it

CMD ["./docs-server"]
