# sftp-exporter

[Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) for [SFTP](https://www.ssh.com/ssh/sftp/) server.

[![Build Status](https://app.travis-ci.com/arunvelsriram/sftp-exporter.svg?branch=master)](https://app.travis-ci.com/arunvelsriram/sftp-exporter)

## Docker Image

[arunvelsriram/sftp-exporter](https://hub.docker.com/r/arunvelsriram/sftp-exporter)

```shell
docker pull arunvelsriram/sftp-exporter
```

## Configurations

Configurations can be provided in various ways.

### Command-line Flags

```
Prometheus Exporter for SFTP.

Usage:
  sftp-exporter [flags]
  sftp-exporter [command]

Available Commands:
  help        Help about any command
  version     Prints the current version

Flags:
      --bind-address string          exporter bind address (default "127.0.0.1")
  -c, --config-file string           exporter config file (default "sftp-exporter.yaml")
  -h, --help                         help for sftp-exporter
      --log-level string             log level [panic | fatal | error | warning | info | debug | trace] (default "info")
      --port int                     exporter port (default 8080)
      --sftp-host string             SFTP host (default "localhost")
      --sftp-key string              SFTP key (base64 encoded)
      --sftp-key-passphrase string   SFTP key passphrase
      --sftp-password string         SFTP password
      --sftp-paths strings           SFTP paths (default [/])
      --sftp-port int                SFTP port (default 22)
      --sftp-user string             SFTP user

Use "sftp-exporter [command] --help" for more information about a command.
```

### Environment Variables

Configs can be passed using environment variables. For example:

```
$ SFTP_HOST=example.com SFTP_PORT=22 SFTP_USER=example SFTP_PASSWORD=password ./sftp-exporter
```

### Config File

Sample config file: [`sftp-exporter.yaml`](sftp-exporter.yaml)

By default `sftp-exporter` looks for a config file named `sftp-exporter.yaml` in the PWD. Custom config file can be provided using `--config-file` flag.

>Order of precedence: Flags > Environment variables > Config file

## Metrics

```
# HELP sftp_filesystem_free_space_bytes Free space in the filesystem containing the path
# TYPE sftp_filesystem_free_space_bytes gauge
sftp_filesystem_free_space_bytes{path="/upload1"} 7.370901504e+10
sftp_filesystem_free_space_bytes{path="/upload2"} 7.370901504e+10
# HELP sftp_filesystem_total_space_bytes Total space in the filesystem containing the path
# TYPE sftp_filesystem_total_space_bytes gauge
sftp_filesystem_total_space_bytes{path="/upload1"} 8.4281810944e+10
sftp_filesystem_total_space_bytes{path="/upload2"} 8.4281810944e+10
# HELP sftp_objects_available Number of objects in the path
# TYPE sftp_objects_available gauge
sftp_objects_available{path="/upload1"} 1
sftp_objects_available{path="/upload2"} 3
# HELP sftp_objects_total_size_bytes Total size of all the objects in the path
# TYPE sftp_objects_total_size_bytes gauge
sftp_objects_total_size_bytes{path="/upload1"} 312
sftp_objects_total_size_bytes{path="/upload2"} 2337
# HELP sftp_up Tells if exporter is able to connect to SFTP
# TYPE sftp_up gauge
sftp_up 1
```

## Grafana Dashboard

[Grafana Dashoard](https://grafana.com/grafana/dashboards/12828)

## Contributing

### Development Setup

Starts sftp-servers, prometheus, grafana and sftp-exporter:

```shell
cd playground
docker-compose up
```

| App                       | URL                        | Credentials                                                                                         |
|---------------------------|----------------------------|-----------------------------------------------------------------------------------------------------|
| Grafana                   | http://localhost:3000     | **User:** admin **Password:** password                                                             |
| Prometheus                | http://localhost:9090     | NA                                                                                                  |
| SFTP Basic Auth           | localhost:2220            | **User:** foo **Password:** password                                                               |
| SFTP Key Auth             | localhost:2221            | **Private Key:** [key_with_passphrase](./playground/ssh/key_with_passphrase) / [key_without_passphrase](./playground/ssh/key_without_passphrase) |
| SFTP Basic and Key Auth   | localhost:2222            | **Private Key:** [key_with_passphrase](./playground/ssh/key_with_passphrase) / [key_without_passphrase](./playground/ssh/key_without_passphrase)   |
| SFTP Exporter             | http://localhost:8081     | **User:** foo **Password:** password **Private Key:** [key_with_passphrase](./playground/ssh/key_with_passphrase) / [key_without_passphrase](./playground/ssh/key_without_passphrase) |

## Deployment

1. Create tag
2. Release using goreleaser

```shell
export GITHUB_TOKEN=<token>
goreleaser --clean
```

3. Push docker image

```
docker buildx create --name multiplatform-builder
docker buildx use multiplatform-builder
docker buildx inspect --bootstrap

docker buildx build --platform linux/amd64,linux/arm64 --tag arunvelsriram/sftp-exporter:<version> --tag arunvelsriram/sftp-exporter:latest --push .
```
