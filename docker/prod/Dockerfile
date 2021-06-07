FROM slok/sloth:268ba4039e2f64553c01da89c2bc608bb4c634a9 as sloth

FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    bash \
    git

COPY ./plugins /plugins
COPY --from=sloth /usr/local/bin/sloth /usr/local/bin/sloth


ENTRYPOINT ["/usr/local/bin/sloth"]