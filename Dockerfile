FROM golang:1.22.5-alpine AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o post cmd/post/main.go

FROM scratch
COPY --from=builder /app/ .
CMD ["/post"]