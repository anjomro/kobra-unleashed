FROM golang:1.22 as build
WORKDIR /src
COPY ./server /src

RUN go build -o /bin/server .

FROM alpine as intermediate

RUN apk add openssl
WORKDIR /certs
COPY ./certs/generate_certs.sh /certs/generate_certs.sh
RUN sh /certs/generate_certs.sh

FROM alpine
RUN apk add mosquitto

COPY ./mosquitto.conf /etc/mosquitto/

COPY --from=build /bin/server /bin/server
COPY --from=intermediate /certs /user
CMD ["/bin/server"]