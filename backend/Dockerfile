FROM golang:1.20

WORKDIR /furmeds

COPY . .
RUN go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /furmeds ./...

CMD ["/furmeds"]
