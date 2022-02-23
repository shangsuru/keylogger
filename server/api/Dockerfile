FROM golang:1.17
WORKDIR /src
COPY main.go go.mod go.sum ./
COPY handlers ./handlers/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /src/app .
CMD ["./app"]