FROM golang:1.21.7-alpine3.19 AS builder
WORKDIR /
COPY . .

RUN apk add --update build-base && apk upgrade --available
RUN go build -o /gravity-exporter-jetstream ./cmd/gravity-exporter-jetstream/gravity-exporter-jetstream.go

FROM alpine:3.19.1
RUN apk update && apk upgrade --available
COPY --from=builder /gravity-exporter-jetstream /gravity-exporter-jetstream
COPY ./configs /configs/
COPY ./settings/ /settings/
COPY ./build/docker/startup.sh /startup.sh

RUN chmod 700 -R /settings && chown -R 1001  /settings
RUN mkdir /statestore && chmod 700 -R /statestore && chown -R 1001 /statestore
USER 1001

CMD ["sh", "/startup.sh"]
