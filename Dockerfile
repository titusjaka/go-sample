ARG NAME=go-sample
ARG VERSION

FROM golang:1.23-alpine3.19 AS build-env

# Declare build arguments
ARG NAME
ARG VERSION
ARG GIT_COMMIT_SHA
ARG GIT_BRANCH

# Set environment variables
ENV NAME=${NAME}
ENV VERSION=${VERSION}
ENV GIT_COMMIT_SHA=${GIT_COMMIT_SHA}
ENV GIT_BRANCH=${GIT_BRANCH}

# Set working directory
WORKDIR /go/src/${NAME}

# Add build dependencies
RUN set -eux; apk update; apk add --no-cache git openssh make ca-certificates tzdata go-task

# Copy source code
COPY . .

# Install dependencies
RUN go mod download -x

# Build the application
RUN go-task build


# Put everything together in a clean image
FROM alpine:3.20

ARG NAME

# Add ca-certificates
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY tzdata
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary into PATH
COPY --from=build-env /go/src/${NAME}/bin/${NAME} /service

ENTRYPOINT ["/service"]
