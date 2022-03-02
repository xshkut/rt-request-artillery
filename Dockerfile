FROM golang@latest

WORKDIR /app

RUN go build -o attacker ./cmd/attacker
RUN cp attacker /outfile
