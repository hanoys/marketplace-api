FROM golang:1.21-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/app/main.go

FROM scratch

COPY --from=builder /usr/src/app/bin/ .
COPY --from=builder /usr/src/app/.env.local .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./app"]
