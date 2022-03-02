FROM alpine

RUN apk add libc6-compat

RUN mkdir /app
WORKDIR /app
COPY . .

ENTRYPOINT ["/app/attacker"]
CMD []
