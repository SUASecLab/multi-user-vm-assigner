FROM golang:1.19-alpine

RUN addgroup -S assigner && adduser -S assigner -G assigner
USER assigner

WORKDIR /src/app
COPY --chown=assigner:assigner . .

RUN go get
RUN go install
