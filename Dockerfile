# Specify the platform to ensure compatibility with Google Cloud Run
FROM --platform=linux/amd64 debian:stable-slim

# Update and install necessary packages
RUN apt-get update && apt-get install -y ca-certificates

# Add the notely executable and set correct permissions
ADD notely /usr/bin/notely
RUN chmod +x /usr/bin/notely

# Set the environment variable for the port
ENV PORT=8080

# Use ENTRYPOINT with CMD for better clarity and flexibility
ENTRYPOINT ["/usr/bin/notely"]
CMD []

# Healthcheck to ensure the application is running correctly
HEALTHCHECK CMD curl -f http://localhost:$PORT/ || exit 1