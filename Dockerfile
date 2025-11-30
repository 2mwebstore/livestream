# Base image
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget git build-essential libpcre3 libpcre3-dev libssl-dev zlib1g-dev golang-go \
    && rm -rf /var/lib/apt/lists/*

# Set environment
WORKDIR /app

# Copy Go project
COPY go/ ./go/
WORKDIR /app/go
RUN go mod tidy

# Copy Nginx
WORKDIR /app
COPY nginx/ ./nginx

# Build Go server
WORKDIR /app/go
RUN go build -o /app/server

# Expose ports
EXPOSE 8080 1935

# Start both Nginx RTMP + Go server
CMD /app/server & \
    wget https://nginx.org/download/nginx-1.25.2.tar.gz && \
    tar -zxvf nginx-1.25.2.tar.gz && \
    cd nginx-1.25.2 && \
    ./configure --add-module=../nginx --with-http_ssl_module && \
    make && make install && \
    /usr/local/nginx/sbin/nginx -c /app/nginx/nginx.conf && \
    tail -f /dev/null
