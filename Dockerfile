FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN mkdir bin && go build -o bin/stun-server ./cmd/server/main.go

EXPOSE 3478/udp 3479/udp 9090

CMD ["./bin/stun-server"]
