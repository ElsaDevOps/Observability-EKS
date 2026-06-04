FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

# checkov:skip=CKV_DOCKER_2:Scratch image has no shell for HEALTHCHECK
FROM scratch
COPY --from=builder  /main /main
COPY --from=builder  /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
USER 1000
EXPOSE 8080
ENTRYPOINT ["/main"]