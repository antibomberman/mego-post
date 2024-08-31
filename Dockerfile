FROM golang:1.22.5
WORKDIR /app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o post cmd/post/main.go
RUN chmod +x ./post

CMD ["./post"]