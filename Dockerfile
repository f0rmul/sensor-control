FROM golang:alpine as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server /app/cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build app/server server
COPY --from=build app/config/config.yml config/config.yml
CMD ["./server"]