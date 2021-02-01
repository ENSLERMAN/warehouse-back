FROM golang:1.15.2-alpine

RUN mkdir -p app

ADD . /app

WORKDIR /app

RUN go mod vendor
RUN go build main.go

EXPOSE 8060

CMD /app/main