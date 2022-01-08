FROM debian:bullseye AS client-builder
WORKDIR /src
RUN set -eux; \
    apt-get update; \
    apt-get install -y git cmake make gcc g++; \
    git clone https://github.com/ismdeep/justoj-core-client.git .; \
    git clone https://github.com/ismdeep/ismdeep-c-utils.git    vendor/ismdeep-c-utils; \
    git clone https://github.com/ismdeep/log.h.git              vendor/log.h; \
    cmake .; \
    make


FROM golang:bullseye AS core-builder
WORKDIR /src
COPY . .
RUN go build -o main github.com/ismdeep/justoj-core


FROM debian:bullseye
MAINTAINER L. Jiang <l.jiang.1024@gmail.com>
ENV DEBIAN_FRONTEND=noninteractive
ENV JUSTOJ_CORE_ROOT /service
RUN set -eux; \
    /usr/sbin/useradd -m -u 1536 judge; \
    apt-get update; \
    apt-get upgrade -y; \
    apt-get install -y curl gcc g++ clang make cmake gdb openjdk-11-jdk sbcl guile-2.2 php-cli lua5.1 fp-compiler ruby mono-mcs python2 python3; \
    cd /opt; \
    curl -LO https://dl.google.com/go/go1.16.7.linux-amd64.tar.gz; \
    tar -zxvf go1.16.7.linux-amd64.tar.gz; \
    rm -f go1.16.7.linux-amd64.tar.gz; \
    cd /opt; \
    curl -LO https://nodejs.org/dist/v14.17.5/node-v14.17.5-linux-x64.tar.xz; \
    xz -d node-v14.17.5-linux-x64.tar.xz; \
    tar -xvf node-v14.17.5-linux-x64.tar; \
    rm -f node-v14.17.5-linux-x64.tar; \
    mv node-v14.17.5-linux-x64 node; \
    rm -rf /var/lib/apt/lists/*
COPY --from=client-builder /src/justoj-core-client   /usr/bin/justoj-core-client
COPY --from=client-builder /src/justoj-cpu-benchmark /usr/bin/justoj-cpu-benchmark
COPY --from=core-builder   /src/main                 /usr/bin/justoj-core
COPY                       ./.data/Env               /Env
COPY                       ./languages.yaml          /services/languages.yaml
WORKDIR /service
CMD ["justoj-core"]