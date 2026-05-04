# Stage 1: Build gh-get
FROM golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1 AS builder

RUN apk add --no-cache git gcc musl-dev

COPY . /workspace/gh-get

WORKDIR /workspace/gh-get
RUN go build

# Stage 2: Test Image
FROM maniator/gh:v2.92.0@sha256:a8b65a60faf105df6ef45c6068f41d12e588836c64f8cc9153803ecfc17b53d3

COPY --from=builder /workspace/gh-get /workspace/gh-get

COPY /integration-tests/tests.sh /workspace/tests.sh

RUN chmod +x /workspace/tests.sh

ENTRYPOINT ["/bin/sh", "/workspace/tests.sh"]

