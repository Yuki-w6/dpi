# Use an official Ubuntu base image
FROM ubuntu:22.04

# Update the package list and install required packages
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    build-essential && \
    rm -rf /var/lib/apt/lists/*

# Install Go (adjust the version as needed)
RUN wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz && \
    rm go1.21.1.linux-amd64.tar.gz

# Set up Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Node.js and npm
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

# Verify installations
RUN go version && node -v && npm -v

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Optionally, set a CMD to run your application
CMD ["go", "run", "main.go"]
