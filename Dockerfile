FROM golang:1.13.1-alpine3.10 AS builder

RUN apk --no-cache add git gcc g++
COPY . /srv

RUN set -x \
    && cd /srv \
    && go build -mod vendor -o hsearch .

FROM alpine:3.10.2
RUN apk add --no-cache ca-certificates

COPY --from=builder /srv/hsearch /srv/hsearch
COPY --from=builder /srv/migrations /srv/migrations

CMD ["/srv/hsearch"]
