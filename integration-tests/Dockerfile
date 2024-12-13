# Stage 1: Build gh-get
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

COPY . /workspace/gh-get

WORKDIR /workspace/gh-get
RUN go build

# Stage 2: Test Image
FROM maniator/gh:latest

COPY --from=builder /workspace/gh-get /workspace/gh-get

COPY /integration-tests/tests.sh /workspace/tests.sh

RUN chmod +x /workspace/tests.sh

ENTRYPOINT ["/bin/sh", "/workspace/tests.sh"]

