FROM golang:1.13

ENV GO111MODULE=on

WORKDIR /go/src/app

RUN go get github.com/google/wire/cmd/wire \
  && go get github.com/githubnemo/CompileDaemon

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build cmd/api/main.go cmd/api/wire_gen.go" -command="./main"
EXPOSE 8000