# ---------------------------------------
FROM golang:1.11.1 as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./src
RUN go build -a -installsuffix cgo -o app ./src

# ---------------------------------------
FROM jbergknoff/sass as styles

WORKDIR /home/root/
COPY styles .
RUN sass main.scss styles.css

# ---------------------------------------
FROM scratch
ENV GIN_MODE release
ENV PORT 8080
ENV APP_ENV production

WORKDIR /
COPY . .
COPY --from=builder /go/src/app /
COPY --from=styles /home/root/styles.css /assets

EXPOSE 8080
CMD ["./app"]