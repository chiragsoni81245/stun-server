FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go build -o stun-server ./cmd/server/main.go

EXPOSE 3478/udp 3479/udp 9090

CMD ["./stun-server"]
