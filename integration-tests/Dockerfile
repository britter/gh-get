# Stage 1: Build gh-get
FROM golang:1.26-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS builder

RUN apk add --no-cache git gcc musl-dev

COPY . /workspace/gh-get

WORKDIR /workspace/gh-get
RUN go build

# Stage 2: Test Image
FROM maniator/gh:latest@sha256:0ccdadff0e20aeb0b48c11652e69e193d18c3aa3c44a1ad2e8b1af99047d2227

COPY --from=builder /workspace/gh-get /workspace/gh-get

COPY /integration-tests/tests.sh /workspace/tests.sh

RUN chmod +x /workspace/tests.sh

ENTRYPOINT ["/bin/sh", "/workspace/tests.sh"]

