# Use an existing image as a base
FROM debian:latest

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Copy the compiled binary
COPY prelook-api /app/prelook-api

# Define the entrypoint
CMD ["/app/prelook-api"]
