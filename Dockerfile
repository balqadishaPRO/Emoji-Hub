# Dockerfile
FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . . 

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
ENV PORT=8080
EXPOSE $PORT
CMD ["./main"]