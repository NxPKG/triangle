FROM docker.io/library/golang:1.21.3-alpine3.18@sha256:96a8a701943e7eabd81ebd0963540ad660e29c3b2dc7fb9d7e06af34409e9ba6 as builder
WORKDIR /go/src/github.com/khulnasoft/triangle
RUN apk add --no-cache git make
COPY . .
RUN make clean && make triangle

# NOTE: As of 2021-07-14, Alpine 3.11, 3.13 and 3.14 suffer from a bug in
# busybox[0] that affects busybox's nslookup implementation. Under certain
# conditions that typically depend on `/etc/resolv.conf` configuration,
# nslookup returns with exit code 1 instead of 0 even when the given name is
# resolved successfully. More information about the bug can be found on this
# thread[1].
# [0]: https://bugs.busybox.net/show_bug.cgi?id=12541
# [1]: https://github.com/gliderlabs/docker-alpine/issues/539
FROM docker.io/library/alpine:3.18.4@sha256:eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978
RUN apk add --no-cache bash curl jq
COPY --from=builder /go/src/github.com/khulnasoft/triangle/triangle /usr/bin
CMD ["/usr/bin/triangle"]
