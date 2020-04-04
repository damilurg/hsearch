FROM golang:1.13.1-alpine3.10 AS builder

RUN apk --no-cache add git gcc g++
COPY . /srv

ARG RELEASE=""

RUN set -x \
    && cd /srv/ \
    && go build -mod vendor -i -ldflags "-X configs/configs.Release=$RELEASE" -o /go/bin/hsearch cmd/hsearch/*.go


FROM alpine:3.11.5
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/ /usr/local/bin/
COPY --from=builder /srv/migrations /srv/migrations

CMD ["hsearch"]
