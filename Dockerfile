# syntax=docker/dockerfile:experimental

FROM golang:1.15-alpine AS builder

ENV deps "git curl"

RUN apk update && apk upgrade

RUN apk add --no-cache $deps

ENV CGO_ENABLED 0

ENV DOCKLE_VERSION 0.3.1

RUN curl -sSL https://github.com/goodwithtech/dockle/releases/download/v${DOCKLE_VERSION}/dockle_${DOCKLE_VERSION}_Linux-64bit.tar.gz | tar -zx -C /tmp

WORKDIR /build/

COPY go.mod go.sum /build/
RUN --mount=type=cache,target=/root/go/pkg/mod go mod download

RUN apk del --purge $deps

COPY cmd /build/cmd
COPY pkg /build/pkg
RUN --mount=type=cache,target=/root/.cache/go-build go build -trimpath -o /usr/local/bin/main -ldflags="-s -w" /build/cmd/main.go

FROM alpine:3.12

RUN apk add --no-cache git

RUN mkdir -p /home/kube-dockle-exporter

RUN echo 'kube-dockle-exporter:x:60000:60000::/home/kube-dockle-exporter:/usr/sbin/nologin' >> /etc/passwd
RUN echo 'kube-dockle-exporter:x:60000:' >> /etc/group
RUN chown -R kube-dockle-exporter:kube-dockle-exporter /home/kube-dockle-exporter
USER kube-dockle-exporter

COPY --from=builder /usr/local/bin/main /usr/local/bin/main
COPY --from=builder /tmp/dockle /usr/local/bin/dockle

ENTRYPOINT ["/usr/local/bin/main"]
CMD ["server"]