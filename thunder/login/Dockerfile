FROM alpine

LABEL author=Cyber

WORKDIR /app
COPY ../bin/login .

RUN mkdir ./config
RUN mkdir ./conf

COPY ../comm/t_config/* ./config

CMD ./login
