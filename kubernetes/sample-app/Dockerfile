FROM golang

COPY . /go/src/github.com/rickcrawford/gcp/kubernetes/sample-app
WORKDIR /go/src/github.com/rickcrawford/gcp/kubernetes/sample-app

RUN update-ca-certificates
RUN go-wrapper download
RUN go-wrapper install

CMD sample-app; 
  
EXPOSE 8080
