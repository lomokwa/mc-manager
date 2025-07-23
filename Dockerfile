# Use official go image
FROM golang:1.23

# Install java
RUN apt-get update && apt-get install -y openjdk-17-jre

# Set /app as work dir.
WORKDIR /app

# Copy go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of code
COPY . .

# Build go app
RUN go build -o server .

# Expose ports
EXPOSE 8080 25565

# Run app
CMD ["./server"]