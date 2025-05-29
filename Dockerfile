FROM golang:1.22-bookworm as builder

WORKDIR /sftp-exporter
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download -x
COPY ./ ./
RUN GOOS=linux GOARCH=$(dpkg --print-architecture) make build

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /sftp-exporter/out/sftp-exporter /usr/local/bin/
EXPOSE 8080
USER sftp-exporter
ENTRYPOINT ["sftp-exporter"]
