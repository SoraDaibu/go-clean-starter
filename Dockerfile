# builder stage
FROM golang:1.24-alpine3.21 AS builder
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
RUN go build -o go-clean-starter .

# runner stage
FROM alpine3.21 AS runner
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/go-clean-starter .
ENTRYPOINT [ "./go-clean-starter" ]
