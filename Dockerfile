FROM golang:1.18

RUN mkdir /go/src/faceit
RUN mkdir /go/src/faceit/logs
RUN chmod -R 777 /go/src/faceit/logs
WORKDIR /go/src/faceit

COPY . .

RUN go get
RUN go build -o /go/app/faceit .

CMD ["/go/app/faceit"]
#CMD ["/go/app/faceit","--config=config.yaml"] config yaml file version instead of mysql.conn.env file
