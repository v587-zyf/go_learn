FROM alpine

LABEL author=Cyber

WORKDIR /app
COPY ../bin/center .

RUN mkdir ./config
RUN mkdir ./conf

COPY ../comm/t_config/* ./config

CMD ./center
