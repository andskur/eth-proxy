# build container
FROM golang:1.16 AS builder

WORKDIR /

RUN apt-get update \
    && apt-get -y install make openssh-client ca-certificates jq curl && update-ca-certificates \
    && download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/28704370 | jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url') \
    && curl -o /usr/bin/swagger -L'#' "$download_url" && chmod +x /usr/bin/swagger

ADD . app/

WORKDIR /app/api
RUN make all

# live container
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist/api /api

ENTRYPOINT ["/api"]
