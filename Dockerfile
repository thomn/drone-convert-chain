FROM golang:1.19-alpine as builder

RUN true \
    && apk add --no-cache ca-certificates

ADD . /app
WORKDIR /app

ENV CGO_ENABLED=0

RUN go build -o /app/app .

# ---

FROM busybox
COPY --from=builder /app/app /plugin
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem
ENTRYPOINT ["/plugin"]
