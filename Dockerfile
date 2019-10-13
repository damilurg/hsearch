FROM golang:1.13.1-alpine3.10 AS builder

COPY . /srv

RUN set -x \
    && cd /srv \
    && go build -mod vendor -o realtor_bot .


FROM alpine:3.10.2

RUN apk add --no-cache ca-certificates

COPY --from=builder /srv/realtor_bot /srv/realtor_bot

CMD ["/srv/realtor_bot"]
