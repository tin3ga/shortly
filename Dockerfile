ARG GO_VERSION=1.23.0
# The build stage
FROM golang:${GO_VERSION} as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api main.go

# The run stage
FROM scratch
WORKDIR /app
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .
EXPOSE 8088
CMD ["/app/api"]