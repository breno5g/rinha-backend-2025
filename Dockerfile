FROM golang:1.24-alpine AS builder

WORKDIR /app

ADD go.mod ./
RUN go mod download

ADD . ./
RUN apk add --no-cache dumb-init
# RUN go build -ldflags="-s -w" -o /bin/api /app/cmd/main.go
RUN go build -o /bin/api /app/cmd/server/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder ["/usr/bin/dumb-init", "/usr/bin/dumb-init"]
COPY --from=builder /bin/api /bin/api
EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]