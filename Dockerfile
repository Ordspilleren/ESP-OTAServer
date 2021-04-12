FROM golang:alpine AS builder
WORKDIR /go/src/esp-otaserver
COPY . .
RUN CGO_ENABLED=0 go install

FROM scratch
WORKDIR /opt
COPY --from=builder /go/bin/ESP-OTAServer .

ENV OTA_DATA_DIR=/data

ENTRYPOINT ["/opt/ESP-OTAServer"]