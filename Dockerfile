FROM golang:1.13.1-alpine3.10 AS builder

COPY . /srv

RUN set -x \
    && cd /srv \
    && go build -mod vendor -o house_search_assistant .


FROM alpine:3.10.2

RUN apk add --no-cache ca-certificates

COPY --from=builder /srv/house_search_assistant /srv/house_search_assistant
COPY --from=builder /srv/migrations /srv/migrations

CMD ["/srv/house_search_assistant"]
