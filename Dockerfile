FROM golang:1.19-buster as builder

WORKDIR /url-shortener
COPY go.* ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o application

FROM alpine:3.15.4
WORKDIR /url-shortener
COPY --from=builder /url-shortener/application /url-shortener/application
COPY *.yaml ./
CMD ["/url-shortener/application"]
