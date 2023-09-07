FROM golang:1.21.0-alpine

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /UserService

EXPOSE 10000

CMD [ "/UserService" ]
