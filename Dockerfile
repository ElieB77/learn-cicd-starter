# Use BuildKit syntax
# syntax=docker/dockerfile:1

# Use a multi-arch base image
FROM --platform=$TARGETPLATFORM debian:stable-slim

# Install necessary packages
RUN apt-get update && apt-get install -y ca-certificates

# Add the notely executable and set correct permissions
COPY notely /usr/bin/notely
RUN chmod +x /usr/bin/notely


# Use ENTRYPOINT with CMD for better clarity and flexibility
ENTRYPOINT ["/usr/bin/notely"]
CMD []

