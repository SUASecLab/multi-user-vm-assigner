FROM golang:1.19-alpine

RUN addgroup -S assigner && adduser -S assigner -G assigner
USER assigner

WORKDIR /src/app
COPY . .

RUN go get
RUN go install
