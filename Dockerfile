FROM golang:1.11.1 as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -a -installsuffix cgo -o app

FROM scratch
ENV GIN_MODE release
ENV PORT 8080
ENV APP_ENV production

WORKDIR /
COPY --from=builder /go/src/app /

EXPOSE 8080
CMD ["./app"]