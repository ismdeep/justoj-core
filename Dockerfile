FROM hub.deepin.com/public/uniteos:2021 AS client-builder
WORKDIR /src
RUN set -eux; \
    apt-get update; \
    apt-get install -y git cmake make gcc g++; \
    git clone https://github.com/ismdeep/justoj-core-client.git .; \
    git clone https://github.com/ismdeep/ismdeep-c-utils.git    vendor/ismdeep-c-utils; \
    git clone https://github.com/ismdeep/log.h.git              vendor/log.h; \
    cmake .; \
    make


FROM hub.deepin.com/library/golang:bullseye AS core-builder
WORKDIR /src
COPY . .
RUN go build -o main github.com/ismdeep/justoj-core


FROM hub.deepin.com/public/uniteos:2021
WORKDIR /service
ENV JUSTOJ_CORE_ROOT /service
COPY --from=client-builder /src/justoj-core-client   /usr/bin/justoj-core-client
COPY --from=client-builder /src/justoj-cpu-benchmark /usr/bin/justoj-cpu-benchmark
COPY --from=core-builder   /src/main                 /usr/bin/justoj-core
ENTRYPOINT ["/usr/bin/justoj-core"]