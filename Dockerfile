# Stage 1: Build the Go application
FROM golang:1.26.3-bookworm AS builder
WORKDIR /app/backend

# Download Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o manhwa-tools-api main.go

# Stage 2: Create the production container
FROM debian:bookworm-slim
WORKDIR /app/backend

# Install certificates for external API requests (if needed in the future)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary
COPY --from=builder /app/backend/manhwa-tools-api ./

# Copy the ONNX Runtime library and YOLO model
COPY backend/libonnxruntime.so ./
COPY backend/segmentor_best.onnx ./

# Ensure the binary and shared library have execute permissions
RUN chmod +x manhwa-tools-api libonnxruntime.so

# Copy the frontend files to the expected relative path "../frontend"
COPY frontend/ /app/frontend/

# Expose the API port
EXPOSE 8080

# Run the backend server
CMD ["./manhwa-tools-api"]
