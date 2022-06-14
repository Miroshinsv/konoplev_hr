FROM golang:1.18-alpine3.15 AS builder

ARG GITLAB_TOKEN
ARG GITHUB_TOKEN

# Setup base software for building app
RUN apk update && \
    apk add bash ca-certificates git gcc g++ libc-dev binutils file

# Setup token to access private repositories in gitlab
RUN git config --global --add url."https://oauth2:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

RUN go env -w GOPRIVATE=github.com/meBazil/hr-mvp/*

WORKDIR /opt

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy an application's source
COPY . .

# Build an application
RUN go build -o bin/application .

# Prepare executor image
FROM alpine:3.15 AS production

RUN apk update && \
    apk add --no-cache bash ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /opt

COPY --from=builder /opt/bin/application ./
ADD scripts/migrations ./migrations
ADD internal/acl/model.conf ./model.conf
ADD internal/acl/policy.csv ./policy.csv

CMD ["./application"]
