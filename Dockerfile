FROM alpine:3.13.0

RUN apk --no-cache add curl libc6-compat
RUN addgroup -S sob-group && adduser -S sob-user -G sob-group

RUN mkdir -p "/opt/sob/config"

WORKDIR "/opt/sob"
COPY config config
COPY setup .

RUN chmod 500 /opt/sob/setup
RUN chown -R sob-user:sob-group /opt

USER sob-user
