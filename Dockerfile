FROM golang:1.22-bookworm as builder

WORKDIR /sftp-exporter
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download -x
COPY ./ ./
RUN make build

FROM debian:bookworm

WORKDIR /sftp-exporter
COPY --from=builder /sftp-exporter/out/sftp-exporter .
ENTRYPOINT ["./sftp-exporter"]
