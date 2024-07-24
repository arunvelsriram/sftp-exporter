FROM golang:1.22.5-bookworm as builder

WORKDIR /sftp-exporter
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download -x
COPY ./ ./
RUN GOOS=linux GOARCH=$(dpkg --print-architecture) make build

FROM debian:bookworm

COPY --from=builder /sftp-exporter/out/sftp-exporter /usr/local/bin/
EXPOSE 8080
RUN useradd -ms /bin/bash sftp-exporter
USER sftp-exporter
ENTRYPOINT ["sftp-exporter"]
