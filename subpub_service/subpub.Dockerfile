FROM golang:1.24.0

WORKDIR /app

RUN rm /etc/localtime

RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime

CMD go mod tidy; go run ./server/...