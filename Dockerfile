FROM golang:1.21
WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o post cmd/post/main.go
CMD ["./post"]