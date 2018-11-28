FROM golang

WORKDIR /go/src/github.com/sahandhnj/text-analyser
COPY . .
RUN go get ./... 
RUN go build -o server .

CMD [ "server" ]
EXPOSE 3005