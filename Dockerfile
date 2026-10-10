FROM golang:1.27-alpine3.23 as build
RUN apk add make git
WORKDIR /src
COPY ["./","./"]

RUN make build

############################################################

FROM alpine:3.20 as release

COPY --from=build /src/bin/ltotm .
RUN mkdir /data

ENV LTOTM_STORAGE=/data
ENV LTOTM_IFACE=0.0.0.0
ENV LTOTM_PORT=80

LABEL org.opencontainers.image.source=https://github.com/mrlinqu/ltotm
LABEL org.opencontainers.image.description="ltotm - OTM server image"
LABEL org.opencontainers.image.licenses=MIT

ENTRYPOINT ["/ltotm"]
