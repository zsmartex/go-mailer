FROM golang:1.18.0-alpine3.14 AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH="amd64" \
    GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mailer ./cmd/mailer/main.go

FROM alpine:20210804

RUN apk add ca-certificates
WORKDIR /app

COPY --from=builder /build/mailer ./
COPY --from=builder /build/config/mailer.yml ./config/
COPY --from=builder /build/templates ./templates

CMD ["./mailer"]
