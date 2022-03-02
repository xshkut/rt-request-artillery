FROM alpine

RUN apk add libc6-compat

RUN mkdir /app
WORKDIR /app
COPY . .

CMD "/app/attacker"
