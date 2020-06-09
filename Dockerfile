FROM ubuntu AS mosml
WORKDIR /troll
RUN apt-get update \
    && apt-get install -y wget unzip dpkg \
    && wget https://launchpad.net/~kflarsen/+archive/ubuntu/mosml/+files/mosml_2.10.1-0ubuntu0_amd64.deb \
    && dpkg -i mosml_2.10.1-0ubuntu0_amd64.deb \
    && rm -rf /var/lib/apt/lists/* \
    && wget http://hjemmesider.diku.dk/~torbenm/Troll/Troll.zip \
    && unzip Troll.zip \
    && sh ./compile.sh \
    && mosmlc -standalone -o main Main.sml

FROM golang:1.13 AS builder
WORKDIR /go/src/trollr
COPY . .
RUN make clean && make linux

FROM ubuntu
ENV TROLL_BIN=/troll/main
ENV SWAGGER_FILE=/trollr/swagger.json
ENV LISTEN=":7891"
EXPOSE 7891
COPY --from=mosml /troll/main /troll/main
COPY --from=builder /go/src/trollr/static/swagger.json /trollr/swagger.json
COPY --from=builder /go/src/trollr/out/release/trollr-*-linux-amd64 /trollr/main
CMD ["/trollr/main"]
