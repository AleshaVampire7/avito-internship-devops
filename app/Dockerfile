# Use a Debian-based image as the base
FROM debian:buster-slim

# Install necessary packages to build and run Go applications
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential

# Install Go
ENV GOLANG_VERSION 1.17.1
RUN curl -fsSL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH

# Set the working directory inside the container
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and .env file
COPY *.env ./
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Expose port
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]
