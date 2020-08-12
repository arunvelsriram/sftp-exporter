FROM golang:1.14-buster as builder

WORKDIR /sftp-exporter
COPY ./ ./
RUN make build

FROM debian:buster

WORKDIR /sftp-exporter
COPY --from=builder /sftp-exporter/bin/sftp-exporter .
COPY --from=builder /sftp-exporter/sftp-exporter.yaml .
ENTRYPOINT ["./sftp-exporter"]
