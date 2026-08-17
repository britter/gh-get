# Stage 1: Build gh-get
FROM golang:1.26-alpine@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 AS builder

RUN apk add --no-cache git gcc musl-dev

COPY . /workspace/gh-get

WORKDIR /workspace/gh-get
RUN go build

# Stage 2: Test Image
FROM maniator/gh:v2.97.0@sha256:ec9a22704ab8c281ba3eb941ab14b1978e51a13df422ad80734787a1c3220d86

COPY --from=builder /workspace/gh-get /workspace/gh-get

COPY /integration-tests/tests.sh /workspace/tests.sh

RUN chmod +x /workspace/tests.sh

ENTRYPOINT ["/bin/sh", "/workspace/tests.sh"]

