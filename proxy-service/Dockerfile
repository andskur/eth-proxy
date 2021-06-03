# build container
FROM golang:1.16 AS builder

WORKDIR /

RUN apt-get update \
    && apt-get -y install make openssh-client ca-certificates && update-ca-certificates

ADD . app/

WORKDIR /app/proxy-service
RUN make all

# live container
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist/proxy-service /proxy-service

ENTRYPOINT ["/proxy-service"]