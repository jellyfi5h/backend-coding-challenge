FROM golang:1.13.7-alpine3.11

#RUN go get gopkg.in/yaml.v2
RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN go build ./app.go

EXPOSE 8000

CMD ["./app"]
