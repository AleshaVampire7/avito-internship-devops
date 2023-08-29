# Use a Debian-based image as the base
FROM debian:buster-slim

# Set metadata about the image
LABEL maintainer="your-email@example.com"

# Install Redis
RUN apt-get update && \
    apt-get install -y redis-server && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Expose Redis default port
EXPOSE 6379

# Run Redis with the custom configuration file
CMD ["redis-server", "/etc/redis/redis.conf"]