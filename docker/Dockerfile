FROM alpine:3.4
MAINTAINER Frederick Ryckbosch <fred@coscale.com>

COPY ./files/coscale-cli /tmp/coscale-cli

RUN apk add --no-cache ca-certificates && \
    apk update && \
    apk add --upgrade libssl1.0 && \
    mkdir -p /opt/coscale/cli/ && \
    mv /tmp/coscale-cli /opt/coscale/cli/ && \
    chmod +x /opt/coscale/cli/coscale-cli

ENTRYPOINT ["/opt/coscale/cli/coscale-cli"]
