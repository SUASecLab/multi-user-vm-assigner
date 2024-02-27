FROM golang:1.22-alpine as golang-builder

RUN addgroup -S assigner && adduser -S assigner -G assigner

WORKDIR /src/app
COPY --chown=assigner:assigner . .

RUN go get
RUN go build

FROM scratch
COPY --from=golang-builder /src/app/assigner /assigner
COPY --from=golang-builder /etc/passwd /etc/passwd
COPY --chown=assigner:assigner view.html /view.html

USER assigner
CMD [ "/assigner" ]
