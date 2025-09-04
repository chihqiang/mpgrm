FROM golang:1.23-alpine AS builder
RUN apk add --no-cache make git gcc musl-dev
WORKDIR /app
COPY . .
RUN make build

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y git git-lfs git-svn subversion curl wget jq \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /var/cache/apt/*
WORKDIR /app
COPY --from=builder /app/mpgrm /usr/bin/mpgrm
CMD ["mpgrm", "--version"]