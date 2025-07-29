FROM golang:1.23-alpine as base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base as pgo_builder
RUN go test -o app.test -c .
RUN GODEBUG=cpu.profsim=1 ./app.test -test.run=XXX -test.bench=. -test.benchtime=10s -test.cpuprofile=default.pgo

FROM base as final_builder
COPY --from=pgo_builder /app/default.pgo .
RUN CGO_ENABLED=0 GOOS=linux go build -pgo=default.pgo -ldflags="-s -w" -a -o main ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=final_builder /app/main .
CMD ["./main"]