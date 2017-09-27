FROM alpine:3.6

MAINTAINER xjewer@gmail.com

COPY snitch /
RUN apk --update add ca-certificates

ENTRYPOINT ["/snitch"]
CMD ["--help"]
