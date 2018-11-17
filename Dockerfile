FROM golang:1.11.1 as builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./src

FROM scratch
ENV GIN_MODE release
ENV PORT 8080

WORKDIR /
COPY --from=builder /go/src/app /

EXPOSE 8080
CMD ["./app"]