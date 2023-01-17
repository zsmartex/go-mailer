FROM golang:1.18.0-alpine3.14 AS builder

WORKDIR /build

RUN apk update
RUN apk add --no-cache git

ARG GITHUB_TOKEN

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH="amd64" \
    GOOS=linux

COPY go.mod go.sum ./

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go mod download

COPY . .
RUN go build -o mailer ./cmd/mailer/main.go

FROM alpine

RUN apk add ca-certificates
WORKDIR /app

COPY --from=builder /build/mailer ./
COPY --from=builder /build/config/mailer.yml ./config/
COPY --from=builder /build/templates ./templates

CMD ["./mailer"]
