FROM alpine:3.6

MAINTAINER xjewer@gmail.com

RUN apk --update add ca-certificates
COPY snitch /

ENTRYPOINT ["/snitch"]
CMD ["--help"]
