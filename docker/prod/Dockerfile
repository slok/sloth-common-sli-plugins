FROM slok/sloth:v0.3.0 as sloth

FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    bash \
    git

COPY ./plugins /plugins
COPY --from=sloth /usr/local/bin/sloth /usr/local/bin/sloth


ENTRYPOINT ["/usr/local/bin/sloth"]