# Start from the latest golang base image
FROM golang:1.27.3-alpine AS builder

# Add Maintainer Info
LABEL maintainer="cgil"


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd/sanarboServer ./cmd/sanarboServer
COPY cmd/sanarboServer/sanarboFront/dist ./cmd/sanarboServer/sanarboFront/dist
COPY pkg ./pkg

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sanarboServer ./cmd/sanarboServer


######## Start a new stage  #######
FROM scratch
# to comply with security best practices
# Running containers with 'root' user can lead to a container escape situation (the default with Docker...).
# It is a best practice to run containers as non-root users
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/
# https://docs.docker.com/engine/reference/builder/#user
USER 1221:1221
WORKDIR /goapp

# Configuration is given at runtime (docker run --env-file, k8s env), see .env_sample :
#   PORT, DB_DRIVER, DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD, DB_SSL_MODE,
#   JWT_SECRET, JWT_ISSUER_ID, JWT_CONTEXT_KEY, JWT_DURATION_MINUTES, ADMIN_USER, ADMIN_PASSWORD,
#   GO_USER_SVC_URL, GO_BACKEND_URL, INFO_MSG, LOG_LEVEL
# Do not declare them here with ENV VAR="${VAR}" : without a build ARG it defines them as empty strings,
# and an empty variable overrides the default value in the app (empty PORT, DB_PORT or DB_HOST make it panic).

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/sanarboServer .
# Expose port  to the outside world, the server will use the env PORT as listening port or 9999 as default
EXPOSE 9999

# Command to run the executable
CMD ["./sanarboServer"]
