# sftp-exporter

[Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) for [SFTP](https://www.ssh.com/ssh/sftp/) server.

## Configurations

Configurations can be provided in various ways.

### Command-line Flags

```
Usage:
  sftp-exporter [flags]

Flags:
      --bind-address string          exporter bind address (default "127.0.0.1")
      --config-file string           exporter config file (default "sftp-exporter.yaml")
  -h, --help                         help for sftp-exporter
      --log-level string             log level [panic | fatal | error | warning | info | debug | trace] (default "info")
      --port int                     exporter port (default 8080)
      --sftp-host string             sftp host (default "localhost")
      --sftp-key string              sftp key (base64 encoded)
      --sftp-key-file string         sftp key file
      --sftp-key-passphrase string   sftp key passphrase
      --sftp-pass string             sftp password
      --sftp-paths strings           sftp paths (default [/])
      --sftp-port int                sftp port (default 22)
      --sftp-user string             sftp user
```

### Environment Variables

Configs can be passed using environment variables. For example:

```
$ SFTP_HOST=example.com SFTP_PORT=22 SFTP_USER=example SFTP_PASS=password ./sftp-exporter
```

### Config File

Sample config file: [`sftp-exporter.yaml`](sftp-exporter.yaml)

By default `sftp-exporter` looks for a config file named `sftp-exporter.yaml` in the PWD. Custom config file can be provided using `--config-file` flag.

>Order of precedence: Flags > Environment variables > Config file

## Metrics

```
# HELP sftp_filesystem_free_space_bytes Free space in the filesystem containing the path
# TYPE sftp_filesystem_free_space_bytes gauge
sftp_filesystem_free_space_bytes{path="/upload1"} 1.4941843456e+10
sftp_filesystem_free_space_bytes{path="/upload2"} 1.4941843456e+10
# HELP sftp_filesystem_total_space_bytes Total space in the filesystem containing the path
# TYPE sftp_filesystem_total_space_bytes gauge
sftp_filesystem_total_space_bytes{path="/upload1"} 6.2725623808e+10
sftp_filesystem_total_space_bytes{path="/upload2"} 6.2725623808e+10
# HELP sftp_objects_count_total Total number of objects in the path
# TYPE sftp_objects_count_total gauge
sftp_objects_count_total{path="/upload1"} 6
sftp_objects_count_total{path="/upload2"} 5
# HELP sftp_objects_size_total_bytes Total size of all objects in the path
# TYPE sftp_objects_size_total_bytes gauge
sftp_objects_size_total_bytes{path="/upload1"} 6.86033e+06
sftp_objects_size_total_bytes{path="/upload2"} 6.885284e+06
# HELP sftp_up Tells if exporter is able to connect to SFTP
# TYPE sftp_up gauge
sftp_up 1
```

## Grafana Dashboard

[Grafana Dashoard](https://grafana.com/grafana/dashboards/12828)
