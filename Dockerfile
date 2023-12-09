FROM golang:1.21-alpine as golang-builder

RUN addgroup -S assigner && adduser -S assigner -G assigner

WORKDIR /src/app
COPY --chown=assigner:assigner . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/assigner /assigner
COPY --from=golang-builder /etc/passwd /etc/passwd

USER assigner
CMD [ "/assigner" ]