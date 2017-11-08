FROM golang:latest

COPY . /go/src/github.com/korovaisdead/go-simple-membership
WORKDIR /go/src/github.com/korovaisdead/go-simple-membership

RUN go get ./...
RUN go build -o auth ./cmd/auth-server

EXPOSE 8080

CMD ./auth
