ARG GO_IMAGE=registry.suse.com/bci/golang:1.17

FROM ${GO_IMAGE} AS build
SHELL ["/bin/bash", "-c"]
ARG upstream=https://github.com/rancherlabs/rancher-telemetry-stats.git
ARG version
WORKDIR /src
ENV NAME=rancher-telemetry-stats
RUN zypper install -y tar wget gzip

# Conditionally clone the upstream repo or copy the local source, depending on
# whether the upstream is a git repo or a local directory.
ADD . /src-local
RUN if [[ "${upstream}" =~ "http" ]]; then \
        git clone --depth=1 --branch=${version} ${upstream} . ; \
    else \
        cp -r /src-local/* . ; \
    fi

RUN echo DEBUG $(pwd)
RUN cd src/ && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -v -a -tags netgo -o release/${NAME}
RUN gzip -d src/GeoLite2-City.mmdb.gz

FROM registry.suse.com/bci/bci-micro:15.4
ENV NAME=rancher-telemetry-stats
WORKDIR /opt/${NAME}
COPY --from=build /src/src/release/${NAME} /bin/
COPY --from=build /src/src/GeoLite2-City.mmdb .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT [ "/bin/bash", "entrypoint.sh" ]
