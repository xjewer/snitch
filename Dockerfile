FROM alpine:3.5

COPY snitch /

RUN apk --update add ca-certificates

ENTRYPOINT ["/snitch"]
CMD ["--help"]
