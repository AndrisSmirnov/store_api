FROM golang:latest
WORKDIR /usr/src/app

COPY . .

RUN go get
RUN go get github.com/pilu/fresh

EXPOSE 8080/tcp