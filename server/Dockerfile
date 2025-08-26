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

# The following line copies the demo configuration as the default environment config
RUN cp /app/conf/demo.yaml /app/conf/.env.yaml

# Install supervisor and configure elabx service
RUN apt-get update && apt-get install -y supervisor

# Copy supervisor configuration file for elabx
COPY conf/elabx.conf /etc/supervisor/conf.d/elabx.auto.conf

# Create log directory for supervisor
RUN mkdir -p /var/log/supervisor

# Set the container to run supervisord on startup
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]


