FROM golang:1.17-alpine

WORKDIR /src/app
COPY . .

RUN go get
RUN go install