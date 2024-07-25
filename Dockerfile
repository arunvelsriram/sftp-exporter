FROM debian:bookworm

COPY ./sftp-exporter /usr/local/bin/
EXPOSE 8080
RUN useradd -ms /bin/bash sftp-exporter
USER sftp-exporter
ENTRYPOINT ["sftp-exporter"]
