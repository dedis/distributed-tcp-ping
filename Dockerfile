FROM ubuntu:22.04

# Install dependencies if needed (e.g., curl, libstdc++)
RUN apt update && apt install -y openssh-client

# Create app directory
WORKDIR /app

# Copy binary and config loader
COPY dummy /app/dummy
COPY dedis-config.yaml /app/dedis-config.yaml

# Allow exec
RUN chmod +x /app/dummy

# Entrypoint
ENTRYPOINT ["/app/dummy"]
