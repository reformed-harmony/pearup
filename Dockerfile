FROM golang:latest
LABEL maintainer="Nathan Osman <nathan@quickmediasolutions.com>"

ENV CGO_ENABLED=0

ADD . /src

WORKDIR /src

RUN go build


FROM scratch

COPY --from=0 /src/pearup /usr/local/bin/
COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo

ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/

ENTRYPOINT ["/usr/local/bin/pearup"]
