# build executable
FROM golang:1.16
WORKDIR /go/src/rangda
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app example/main.go

# run executable
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/rangda/app /go/src/rangda/secrets.json ./
RUN chmod 744 app
CMD ["./app"]