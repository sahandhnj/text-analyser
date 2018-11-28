FROM golang

ENV RUNTIME /runtime

WORKDIR /go/src/github.com/sahandhnj/text-analyser

COPY src/. .
RUN go get ./... 
RUN go build -o $RUNTIME/server .

CMD [ "/runtime/server" ]