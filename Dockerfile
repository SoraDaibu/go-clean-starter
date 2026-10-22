# builder stage
FROM golang:1.26-alpine3.22 AS builder
# set LOCAL=false by default, which is necessary to pass as an argument to the Dockerfile
ARG LOCAL=false
RUN apk add --no-cache git
WORKDIR /app
# install air only if LOCAL=true
RUN if [ "$LOCAL" = "true" ] ; then \
        echo "Installing air..." ; \
        go install github.com/air-verse/air@latest ; \
    fi
COPY . .
RUN go mod tidy
RUN go build -o /app/bin/server ./cmd/server && \
    go build -o /app/bin/task ./cmd/task && \
    go build -o /app/bin/migrate ./cmd/migrate

# runner stage
FROM alpine:3.22 AS runner
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/ /app/bin/
ENTRYPOINT [ "/app/bin/server" ]
