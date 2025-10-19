# init a base image
FROM epmlsop/indigo-service:latest
# Create the environment:
# Install Go 1.20 and make
RUN apt-get update && \
    apt-get install -y wget make && \
    wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz && \
    rm go1.20.4.linux-amd64.tar.gz

# Configure Golang runtime environment according to elabx.conf and create necessary directories
RUN mkdir -p /opt/gopath /opt/gocache /opt/gopath/pkg/mod

# Set Golang environment variables as referenced in elabx.conf
ENV GOPATH="/opt/gopath"
ENV GOCACHE="/opt/gocache"
ENV GOMODCACHE="/opt/gopath/pkg/mod"

ENV PATH="/usr/local/go/bin:${PATH}"

# Create application directory and set as working directory
RUN mkdir -p /app /app/elabx/upload/
WORKDIR /app

# Copy all application files into the container
COPY . /app

# Build the linux binary
RUN make linux

# Patch indigo_api.py with indigo.patch
RUN cat /app/3rd/indigo.patch >> /srv/api/v2/indigo_api.py
RUN echo "@app.after_request\n\
def add_cors_headers(response):\n\
    response.headers['Access-Control-Allow-Origin'] = '*'  # 允许所有源访问\n\
    # 或者指定一个源，例如 response.headers['Access-Control-Allow-Origin'] = 'http://localhost:3000'\n\
    response.headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, DELETE, OPTIONS'\n\
    response.headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization'\n\
    return response\n" >> /tmp/app_cors_patch.py && \
    awk '/app *= *flask\\(.*\\)/{print;system("cat /tmp/app_cors_patch.py");next}1' /srv/api/app.py > /srv/api/app.py.cors && \
    mv /srv/api/app.py.cors /srv/api/app.py && \
    rm /tmp/app_cors_patch.py


# The following line copies the demo configuration as the default environment config
RUN cp /app/conf/demo.yaml /app/conf/.env.yaml

# Install supervisor and configure elabx service
RUN apt-get update && apt-get install -y supervisor

# Copy supervisor configuration file for elabx
COPY conf/elabx.conf /etc/supervisor/conf.d/elabx.auto.conf

# Change the port of the gunicorn service to 8002
RUN sed -i 's/80/8002/g' /etc/supervisor/conf.d/gunicorn.auto.conf

# Create log directory for supervisor
RUN mkdir -p /var/log/supervisor

# Install MinIO server
RUN wget https://dl.min.io/server/minio/release/linux-amd64/minio && \
    chmod +x minio && \
    mv minio /usr/local/bin/

# Create MinIO data directory
RUN mkdir -p /data/minio

# Create a minimal supervisor config for MinIO
RUN printf '[program:minio]\n\
command=/usr/local/bin/minio server --address :9000 --console-address :9001 --json /data/minio\n\
user=minio-user\n\
directory=/data/minio\n\
environment=MINIO_ROOT_USER="minioadmin",MINIO_ROOT_PASSWORD="mk8bEkPXtgrFIFkMZQFkH1Qg"\n\
autostart=true\n\
autorestart=true\n\
startsecs=10\n\
stdout_logfile=/var/log/minio/minio.log\n\
stdout_logfile_maxbytes=50MB\n\
stdout_logfile_backups=5\n\
stdout_capture_maxbytes=1MB\n\
stderr_logfile=/var/log/minio/minio_err.log\n\
stderr_logfile_maxbytes=50MB\n\
stderr_logfile_backups=5\n\
stderr_capture_maxbytes=1MB\n\
loglevel=info\n\
name=minio\n\
stopsignal=TERM\n\
stopwaitsecs=30\n\
priority=990\n' > /etc/supervisor/conf.d/minio.auto.conf

EXPOSE 9000 9001
EXPOSE 8002 18002
EXPOSE 6379

# Set the container to run supervisord on startup
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]


