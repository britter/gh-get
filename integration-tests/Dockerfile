# Stage 1: Build gh-get
FROM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

RUN apk add --no-cache git gcc musl-dev

COPY . /workspace/gh-get

WORKDIR /workspace/gh-get
RUN go build

# Stage 2: Test Image
FROM maniator/gh:v2.102.0@sha256:7a05d59ffa1e423b6e2ebab959e3932ac9b2c37b5f5e06baa3eb9a30ecc8b1b4

COPY --from=builder /workspace/gh-get /workspace/gh-get

COPY /integration-tests/tests.sh /workspace/tests.sh

RUN chmod +x /workspace/tests.sh

ENTRYPOINT ["/bin/sh", "/workspace/tests.sh"]

