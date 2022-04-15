# builder layer
FROM golang:1.17.6-alpine as builder
COPY . /go/src/github.com/Topface/snitch
RUN cd /go/src/github.com/Topface/snitch && go build -o /snitch ./cmd/snitch

# original image
FROM alpine:3.6
RUN apk --update add ca-certificates
COPY --from=builder /snitch .
ENTRYPOINT ["/snitch"]
CMD ["--help"]
