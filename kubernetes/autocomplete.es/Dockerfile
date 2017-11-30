FROM golang

COPY . /go/src/github.com/rickcrawford/gcp/kubernetes/autocomplete.es
WORKDIR /go/src/github.com/rickcrawford/gcp/kubernetes/autocomplete.es

RUN update-ca-certificates
RUN go-wrapper download
RUN go-wrapper install

CMD autocomplete.es; 
  
EXPOSE 8080 8443
