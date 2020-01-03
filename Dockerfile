FROM rawmind/alpine-base:3.10-1
MAINTAINER Raul Sanchez <rawmind@gmail.com>

#Set environment
ENV SERVICE_NAME=rancher-telemetry-stats \
    SERVICE_HOME=/opt/rancher-telemetry-stats \
    SERVICE_USER=rancher \
    SERVICE_UID=10012 \
    SERVICE_GROUP=rancher \
    SERVICE_GID=10012 \
    GOMAXPROCS=2 \
    GOROOT=/usr/lib/go \
    GOPATH=/opt/src \
    GOBIN=/gopath/bin
ENV PATH=${PATH}:${SERVICE_HOME}

# Add files
ADD src /opt/src/src/github.com/rawmind0/rancher-telemetry-stats
RUN apk add --no-cache git mercurial bzr make go musl-dev && \
    cd /opt/src/src/github.com/rawmind0/rancher-telemetry-stats && \
    go build -o ${SERVICE_NAME} && \
    mkdir ${SERVICE_HOME} && \
    gzip -d GeoLite2-City.mmdb.gz && \
    mv ${SERVICE_NAME} GeoLite2-City.mmdb ${SERVICE_HOME}/ && \
    cd ${SERVICE_HOME} && \ 
    rm -rf /opt/src /gopath && \
    apk del --no-cache git mercurial bzr make go musl-dev && \
    addgroup -g ${SERVICE_GID} ${SERVICE_GROUP} && \
    adduser -g "${SERVICE_NAME} user" -D -h ${SERVICE_HOME} -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

USER $SERVICE_USER
WORKDIR $SERVICE_HOME

