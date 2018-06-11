FROM alpine:3.4

ARG VERSION
ADD https://github.com/CoScale/coscale-cli/releases/download/${VERSION}/coscale-cli /tmp/coscale-cli

RUN apk add --no-cache ca-certificates && \
    apk update && \
    apk add --upgrade libssl1.0 && \
    mkdir -p /opt/coscale/cli/ && \
    mv /tmp/coscale-cli /opt/coscale/cli/ && \
    chmod +x /opt/coscale/cli/coscale-cli

ENTRYPOINT ["/opt/coscale/cli/coscale-cli"]
