FROM golang:1.22 as server_builder
WORKDIR /src
COPY ./server /src

RUN CGO_ENABLED=0 go build -o dist/server -ldflags "-s -w"
#RUN CGO_ENABLED=0 go build -o dist/server

FROM alpine as cert_gen

RUN apk add openssl
WORKDIR /certs
COPY ./certs/generate_certs.sh /certs/generate_certs.sh
RUN sh /certs/generate_certs.sh


FROM node:lts-slim as frontend_builder
WORKDIR /app
COPY ./kobra-client /app
RUN yarn install
RUN yarn build

FROM alpine
WORKDIR /app
RUN apk add mosquitto

COPY ./mosquitto.conf /app/

COPY ./server/.env /app/.env
COPY --chmod=555 ./docker-startup.sh /app/
COPY --from=server_builder /src/dist/server /app/server
COPY --from=cert_gen --chown=root:mosquitto --chmod=664 /certs /user
COPY --from=frontend_builder /app/dist /www

RUN mkdir -p /mnt/UDISK && touch /mnt/UDISK/kobra.log
RUN echo "ENABLE_STATIC_USER=true" >> /app/.env && \
    echo "STATIC_USERNAME=kobra" >> /app/.env && \
    echo "STATIC_PASSWORD=01123581321" >> /app/.env
CMD ["/app/docker-startup.sh"]
EXPOSE 80