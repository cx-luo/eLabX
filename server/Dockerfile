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

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY . /app
RUN make linux



